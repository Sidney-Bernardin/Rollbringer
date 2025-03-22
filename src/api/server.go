package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"rollbringer/src"
	"rollbringer/src/domain/accounts"
	"rollbringer/src/domain/play"
)

type server struct {
	*http.Server

	log    *slog.Logger
	config *src.Config

	accounts accounts.Service
	play     play.Service
}

func NewServer(log *slog.Logger, config *src.Config, accountsSvc accounts.Service, playSvc play.Service) *server {
	svr := &server{
		&http.Server{
			Addr:    config.APIAddr,
			Handler: chi.NewRouter(),
		},
		log, config, accountsSvc, playSvc,
	}

	r := svr.Server.Handler.(chi.Router)
	r.Use(svr.mwLog())
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(os.DirFS("src/api/static"))))

	r.Route("/pages", func(pages chi.Router) {
		pages.With(svr.mwInstance("page-home")).
			Get("/", svr.handlePageHome)
	})

	r.Route("/users", func(users chi.Router) {
		users.With(svr.mwInstance("user-create")).
			Post("/", svr.handleUserCreate)

		users.With(svr.mwInstance("user-get")).
			Get("/{username}", svr.handleUserGet)
	})

	r.Route("/rooms", func(rooms chi.Router) {
		rooms.With(svr.mwInstance("room-create")).
			Post("/", svr.handleRoomCreate)

		rooms.With(svr.mwInstance("room-get")).
			Get("/{room_id}", svr.handleRoomGet)
	})

	return svr
}

func (svr *server) state(r *http.Request) map[string]any {
	state, ok := r.Context().Value("state").(map[string]any)
	if !ok {
		state = map[string]any{}
		*r = *r.WithContext(context.WithValue(r.Context(), "state", state))
	}
	return state
}

func (api *server) logServerError(ctx context.Context, err error) {
	api.log.Log(ctx, src.LevelError,
		"Internal Server Error", "err", err.Error())
}
