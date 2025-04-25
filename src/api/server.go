package api

import (
	"context"
	"embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"

	"rollbringer/src"
	"rollbringer/src/services"
	"rollbringer/src/services/accounts"
	"rollbringer/src/services/play"
)

//go:embed static
var static embed.FS

var (
	externalErrorTypeInvalidProvider src.ExternalErrorType = "invalid_provider"
)

type server struct {
	*http.Server

	log    *slog.Logger
	config *src.Config
	broker services.Broker

	accounts         accounts.Service
	accountsDatabase accounts.BasicDatabase
	google           accounts.Google
	spotify          accounts.Spotify

	play         play.Service
	playDatabase play.Database
}

func NewServer(
	log *slog.Logger,
	config *src.Config,
	broker services.Broker,
	accountsSvc accounts.Service,
	accountsDB accounts.BasicDatabase,
	google accounts.Google,
	spotify accounts.Spotify,
	playSvc play.Service,
	playDB play.Database,
) *server {
	svr := &server{
		&http.Server{
			Addr: config.APIAddr,
		},
		log, config, broker,
		accountsSvc, accountsDB, google, spotify,
		playSvc, playDB,
	}

	r := http.NewServeMux()
	r.Handle("GET /static/", http.FileServerFS(static))

	r.Handle("GET /login/{provider}", svr.handleOAuthConsent())
	r.Handle("GET /login/{provider}/callback", svr.handleOAuthCallback())

	r.Handle("POST /rooms", mw(svr.mwAuth(true, true, ""))(svr.handleRoomCreate()))

	r.Handle("GET /", mw(svr.mwAuth(false, false, "/"))(svr.handlePageHome()))
	r.Handle("GET /play", mw(svr.mwAuth(true, false, "/"))(svr.handlePagePlay()))
	r.Handle("GET /play/ws", mw(svr.mwAuth(true, false, ""))(svr.handlePagePlayWebSocket()))

	svr.Server.Handler = mw(svr.mwLog)(r)
	return svr
}

func (api *server) logServerError(ctx context.Context, err error) {
	api.log.Log(ctx, src.LevelError, "Internal Server Error", "err", err.Error())
}

func decodeEvent(bEvent []byte, events map[string]any) any {
	var head struct {
		Operation string `json:"operation"`
	}

	if err := json.Unmarshal(bEvent, &head); err != nil {
		return nil
	}

	event, ok := events[head.Operation]
	if !ok {
		return nil
	}
	event = reflect.New(reflect.ValueOf(event).Elem().Type()).Interface()

	if err := json.Unmarshal(bEvent, event); err != nil {
		return nil
	}

	return event
}
