package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"rollbringer/src"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/domain/play"
)

var (
	externalErrorTypeUnauthorized    src.ExternalErrorType = "unauthorized"
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

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(os.DirFS("src/api/static"))))

	r.Handle("GET /login/{provider}", svr.handleOAuthConsent())
	r.Handle("GET /login/{provider}/callback", svr.handleOAuthCallback())

	r.Handle("POST /rooms", svr.handleRoomCreate())
	r.Handle("GET /rooms/{room_id}", svr.handleRoomGet())

	r.Handle("/", svr.handlePageHome())

	svr.Handler = mw(svr.mwLog)(r)
	return svr
}

func (api *server) logServerError(ctx context.Context, err error) {
	api.log.Log(ctx, src.LevelError, "Internal Server Error", "err", err.Error())
}
