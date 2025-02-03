package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/repositories/nats"

	"github.com/dpotapov/slogpfx"
	"github.com/go-chi/chi/v5"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

type feature struct {
	handler http.Handler
	service domain.IService
}

var (
	registeredFeatures = map[string]func() error{}
	features           = map[string]feature{}

	config *domain.Config
	logger *slog.Logger

	httpServer *http.Server
)

func main() {
	ctx := context.Background()

	// Create dependencies.
	if err := dependencies(); err != nil {
		domain.HandleError(ctx, logger, domain.SlogLevelFatal, domain.Wrap(err, "cannot create dependencies", nil))
		return
	}

	// Create features.
	for name, createFeature := range registeredFeatures {
		if err := createFeature(); err != nil {
			domain.HandleError(ctx, logger, domain.SlogLevelFatal, domain.Wrap(err, "cannot create feature", map[string]any{"feature": name}))
			return
		}
	}

	// Run!
	err := run(ctx)
	domain.HandleError(ctx, logger, domain.SlogLevelFatal, domain.Wrap(err, "cannot run", nil))
	logger.InfoContext(ctx, "Gracefully shutting down")

	// Shutdown HTTP server.
	if err := httpServer.Shutdown(ctx); err != nil {
		domain.HandleError(ctx, logger, domain.SlogLevelFatal, domain.Wrap(err, "cannot shutdown HTTP server", nil))
	}

	// Shutdown services.
	errChan := make(chan error, 1)
	for name, feat := range features {
		go func() {
			err := feat.service.Shutdown(ctx)
			errChan <- domain.Wrap(err, "cannot shutdown service", map[string]any{"feature": name})
		}()
	}

	// Wait for services to shutdown.
	for range len(features) {
		domain.HandleError(ctx, logger, domain.SlogLevelFatal, <-errChan)
	}

	logger.InfoContext(ctx, "Goodbye, World!")
}

func dependencies() (err error) {
	logger = slog.Default()

	// Create config.
	config, err = domain.GetConfig()
	if err != nil {
		return domain.Wrap(err, "cannot get configuration", nil)
	}

	// Create log handler
	var h slog.Handler
	if config.LogPretty {

		// Make logs pretty and human readable.
		h = tint.NewHandler(os.Stderr, &tint.Options{
			Level:   domain.SlogLevelTrace,
			NoColor: config.LogPrettyColor && !isatty.IsTerminal(os.Stderr.Fd()),
		})

		// Prominently displays "dependency" attributes.
		h = slogpfx.NewHandler(h, &slogpfx.HandlerOptions{
			PrefixKeys: []string{"dependency"},
		})
	} else {

		// Converts logs to JSON.
		h = slog.NewJSONHandler(os.Stderr, nil)
	}

	// Create logger.
	logger = slog.New(h)

	if config.NATSEmbeddedServer {

		// Create embedded NATS server.
		if err := nats.CreateEmbeddedServer(config, logger.With("dependency", "nats-embedded-server")); err != nil {
			return domain.Wrap(err, "cannot create embedded NATS server", nil)
		}
	}

	return nil
}

func run(ctx context.Context) error {
	ctx, sigCancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer sigCancel()

	ctx, causeCancel := context.WithCancelCause(ctx)
	defer causeCancel(nil)

	// Create HTTP server.
	httpServer = &http.Server{
		Addr:    config.Addr,
		Handler: chi.NewRouter(),
	}

	for name, feat := range features {

		// Add the feature's handler to the HTTP server.
		httpServer.Handler.(chi.Router).Mount("/"+name, feat.handler)

		// Run the feature's service.
		go func() {
			if err := feat.service.Run(ctx); err != nil {
				causeCancel(domain.Wrap(err, "cannot run service", map[string]any{"feature": name}))
			}
		}()
	}

	// Run HTTP server.
	go func() {
		logger.InfoContext(ctx, "Running HTTP server", "addr", config.Addr)
		err := httpServer.ListenAndServe()
		causeCancel(domain.Wrap(err, "cannot run HTTP server", nil))
	}()

	<-ctx.Done()
	sigCancel()

	// Check the cause of the context's cancellation.
	if err := context.Cause(ctx); err != nil && err != context.Canceled {
		return err
	}

	return nil
}
