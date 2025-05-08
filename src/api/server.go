package api

import (
	"embed"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/src"
	"rollbringer/src/api/views"
	"rollbringer/src/domain"
	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/domain/services/play"
)

//go:embed static
var static embed.FS

var (
	externalErrorTypeInvalidProvider domain.ExternalErrorType = "invalid-provider"
)

type server struct {
	*http.Server

	log    *slog.Logger
	config *src.Config
	broker domain.PublicBroker

	accounts         accounts.Service
	accountsDatabase accounts.DatabaseCommon
	google           accounts.Google
	spotify          accounts.Spotify

	play         play.Service
	playBroker   play.BrokerCommon
	playDatabase play.Database
}

func NewServer(
	log *slog.Logger,
	config *src.Config,
	broker domain.PublicBroker,

	accountsSvc accounts.Service,
	accountsDB accounts.DatabaseCommon,
	google accounts.Google,
	spotify accounts.Spotify,

	playSvc play.Service,
	playBroker play.BrokerCommon,
	playDB play.Database,
) *server {
	svr := &server{
		&http.Server{
			Addr: config.APIAddr,
		},
		log, config, broker,
		accountsSvc, accountsDB, google, spotify,
		playSvc, playBroker, playDB,
	}

	r := http.NewServeMux()
	r.Handle("GET /static/", http.FileServerFS(static))

	r.Handle("GET /login/{provider}", svr.handleOAuthConsent())
	r.Handle("GET /login/{provider}/callback", svr.handleOAuthCallback())

	r.Handle("POST /rooms", mw(svr.mwAuth(true, true, ""))(svr.handleRoomCreate()))
	r.Handle("GET /rooms/{room_id}/ws", mw(svr.mwAuth(true, false, ""))(svr.handleRoomWebSocket()))

	r.Handle("GET /", mw(svr.mwAuth(false, false, "/"))(svr.handlePageHome()))
	r.Handle("GET /play", mw(svr.mwAuth(true, false, "/"))(svr.handlePagePlay()))

	svr.Server.Handler = mw(svr.mwLog)(r)
	return svr
}

func (svr *server) subChat(conn *websocket.Conn, r *http.Request, roomID uuid.UUID) {
	var ctx = r.Context()

	err := svr.playBroker.SubChat(ctx, roomID, func(event *play.EventChat) {

		// Get the chat event's author.
		author, err := svr.accounts.GetUserByUserID(ctx, uuid.MustParse(event.AuthorID))
		if err != nil {
			svr.err(conn, r, errors.Wrap(err, "cannot get author by author-ID"))
			return
		}

		svr.respond(conn, r, 0, views.ChatMessage(author, event.Message))
	})

	svr.err(conn, r, errors.Wrap(err, "cannot subscribe to chat"))
	conn.Close()
}

func (svr *server) subRoom(conn *websocket.Conn, r *http.Request, roomID uuid.UUID) {
	var ctx = r.Context()

	err := svr.playBroker.SubRoom(ctx, roomID, func(event any) {
		switch e := event.(type) {
		case *domain.EventRoomJoined:
			svr.respond(conn, r, 0, views.NewUserBubble(e))
		}
	})

	svr.err(conn, r, errors.Wrap(err, "cannot subscribe to chat"))
	conn.Close()
}

func (svr *server) subUser(conn *websocket.Conn, r *http.Request, userID uuid.UUID) {
	var ctx = r.Context()

	err := svr.playBroker.SubUser(ctx, userID, func(event any) {
		switch e := event.(type) {
		case *domain.EventNewBoard:
			svr.respond(conn, r, 0, views.NewBoardCard(e))
		}
	})

	svr.err(conn, r, errors.Wrap(err, "cannot subscribe to chat"))
	conn.Close()
}

func (svr *server) respond(w io.Writer, r *http.Request, statusCode int, res any) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.WriteHeader(statusCode)
	}

	var err error
	switch res := res.(type) {
	case []byte:
		_, err = w.Write(res)

	case templ.Component:
		err = res.Render(r.Context(), w)
		err = errors.Wrap(err, "cannot render Templ response")

	default:
		err = json.NewEncoder(w).Encode(res)
		err = errors.Wrap(err, "cannot marshal response")
	}

	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
		svr.log.Log(r.Context(), src.LevelError, "Cannot write response", "err", err.Error())
	}
}

var errCodes = map[domain.ExternalErrorType]int{
	externalErrorTypeInvalidProvider: http.StatusBadRequest,

	domain.ExternalErrorTypeInternalError: http.StatusInternalServerError,
	domain.ExternalErrorTypeUnauthorized:  http.StatusUnauthorized,
	domain.ExternalErrorTypeInvalidUUID:   http.StatusUnprocessableEntity,

	accounts.ExternalErrorTypeInvalidUsername:       http.StatusBadRequest,
	accounts.ExternalErrorTypeProviderNotLinked:     http.StatusBadRequest,
	accounts.ExternalErrorTypeProviderAlreadyLinked: http.StatusConflict,

	play.ExternalErrorTypeInvalidRoomName: http.StatusBadRequest,
}

func (svr *server) err(w io.Writer, r *http.Request, err error) {
	if err == nil {
		return
	}

	// If the error isn't a domain.ExternalError, treat it as a server error.
	var externalErr *domain.ExternalError
	if !errors.As(err, &externalErr) {
		svr.log.Log(r.Context(), src.LevelError, "Internal server error", "err", err.Error())
		externalErr = &domain.ExternalError{Type: domain.ExternalErrorTypeInternalError}
	}

	switch w.(type) {
	case *websocket.Conn:
		svr.respond(w, r, 0, &views.WebSocketResponse{
			Operation: "error",
			Payload:   externalErr,
		})

	case http.ResponseWriter:
		svr.respond(w, r, errCodes[externalErr.Type], externalErr)
	}
}
