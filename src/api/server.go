package api

import (
	"context"
	"embed"
	"log/slog"
	"net/http"

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

func (svr *server) roomCallback(conn *websocket.Conn, r *http.Request) func(any) {
	return func(event any) {
		switch e := event.(type) {
		case *domain.EventRoomJoined:
			svr.respond(conn, r, 0, views.NewUserBubble(e))
		}
	}
}

func (svr *server) chatCallback(ctx context.Context, conn *websocket.Conn, r *http.Request) func(*play.EventChat) {
	return func(event *play.EventChat) {

		// Get the chat message's author.
		author, err := svr.accounts.GetUserByUserID(ctx, uuid.MustParse(event.AuthorID))
		if err != nil {
			if errors.Is(err, domain.ErrEntityNotFound) {
				return
			}

			svr.err(conn, r, errors.Wrap(err, "cannot get chat message's author"))
			return
		}

		svr.respond(conn, r, 0, views.ChatMessage(author, event.Message))
	}
}

func (svr *server) userCallback(conn *websocket.Conn, r *http.Request) func(any) {
	return func(event any) {
		switch e := event.(type) {
		case *domain.EventNewBoard:
			svr.respond(conn, r, 0, views.NewBoardCard(e))
		}
	}
}
