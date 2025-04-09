package api

import (
	"context"
	"embed"
	"log/slog"
	"net/http"

	"rollbringer/src"
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

	accounts   accounts.Service
	accountsDB accounts.DatabaseQueries
	google     accounts.Google
	spotify    accounts.Spotify

	play   play.Service
	playDB play.DatabaseQueries
}

func NewServer(
	log *slog.Logger,
	config *src.Config,
	accountsSvc accounts.Service,
	accountsDB accounts.DatabaseQueries,
	google accounts.Google,
	spotify accounts.Spotify,
	playSvc play.Service,
	playDB play.DatabaseQueries,
) *server {
	svr := &server{
		&http.Server{
			Addr: config.APIAddr,
		},
		log, config,
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

	svr.Server.Handler = mw(svr.mwLog)(r)
	return svr
}

func (api *server) logServerError(ctx context.Context, err error) {
	api.log.Log(ctx, src.LevelError, "Internal Server Error", "err", err.Error())
}
