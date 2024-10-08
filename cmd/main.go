package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/internal/config"
	"rollbringer/internal/repositories/database"
	databases "rollbringer/internal/repositories/database"
	"rollbringer/internal/services"
)

type globalDependencies struct {
	cfg    *config.Config
	logger *slog.Logger

	dbRepo *databases.Database
}

var features = map[string]func(globalDependencies) (http.Handler, services.BaseServicer, error){}

func main() {

	// Create logger.
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))

	// Create configuration.
	cfg, err := config.New()
	if err != nil {
		logger.Error("Cannot create configuration", "err", err.Error())
		return
	}

	// Create a Database repository.
	dbRepo, err := database.New(cfg, logger)
	if err != nil {
		logger.Error("Cannot create database repository", "err", err.Error())
		return
	}

	var (
		router   = chi.NewRouter()
		services = map[string]services.BaseServicer{}
	)

	for name, initFeature := range features {
		handler, service, err := initFeature(globalDependencies{
			cfg:    cfg,
			logger: logger,
			dbRepo: dbRepo,
		})

		if err != nil {
			logger.Error("Cannot create "+name+" feature", "err", err.Error())
			return
		}

		router.Mount("/"+name, handler)
		services[name] = service
	}

	router.Handle("/static/*", handleStaticDir())

	run(cfg, logger, services, &http.Server{
		Handler: router,
		Addr:    ":" + cfg.Port,
	})
}

func run(cfg *config.Config, logger *slog.Logger, services map[string]services.BaseServicer, svr *http.Server) {
	errChan := make(chan error)

	// Run the services.
	for name, svc := range services {
		go func() {
			logger.Info("Running " + name + " service")
			err := svc.Listen()
			errChan <- errors.Wrapf(err, "cannot run %s service", name)
		}()
	}

	// Run the HTTP server.
	go func() {
		logger.Info("Listening on " + svr.Addr)
		err := svr.ListenAndServe()
		errChan <- errors.Wrap(err, "cannot run HTTP server")
	}()

	// Create a channel for Interrupt signals.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Wait...
wait:
	select {
	case sig := <-signalChan:
		logger.Info("Signal interruption", "signal", sig)
	case err := <-errChan:
		if err == nil {
			goto wait
		}
		logger.Error("Fatal error", "err", err.Error())
	}

	shutdown(cfg, logger, services, svr)
}

func shutdown(cfg *config.Config, logger *slog.Logger, services map[string]services.BaseServicer, svr *http.Server) {
	logger.Info("Gracefully shutting down")
	defer logger.Info("Goodbye, World!")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Gracefully shutdown the HTTP server.
	if err := svr.Shutdown(ctx); err != nil {
		logger.Error("Cannot gracefully shutdown HTTP server", "err", err.Error())
	}

	// Gracefully shutdown the services.
	wg := &sync.WaitGroup{}
	for name, svc := range services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := svc.Shutdown(); err != nil {
				logger.Error("Cannot gracefully shutdown "+name+" service", "err", err.Error())
			}
		}()
	}
	wg.Wait()
}
