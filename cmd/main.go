package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"rollbringer/internal/config"
)

var serviceHandlers = map[string]func(*config.Config, *slog.Logger) (http.Handler, error){}

func main() {

	cfg, err := config.New()
	if err != nil {
		slog.Error("Cannot create configuration", "err", err.Error())
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))

	router := chi.NewRouter()
	router.Handle("/static/*", handleStaticDir())

	for pattern, fn := range serviceHandlers {
		svcHandler, err := fn(cfg, logger)
		if err != nil {
			logger.Error("Cannot create service-handler", "err", err.Error())
			return
		}
		router.Mount(pattern, svcHandler)
	}

	err = http.ListenAndServe(cfg.Address, router)
	logger.Error("Server stopped", "err", err.Error())
}
