//go:build all || users
// +build all users

package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"rollbringer/internal/config"
	handler "rollbringer/internal/handlers/users"
	database "rollbringer/internal/repositories/databases/users"
	"rollbringer/internal/repositories/oauth"
	"rollbringer/internal/repositories/pubsub"
	service "rollbringer/internal/services/users"
)

func init() {
	serviceHandlers["/users"] = func(ctx context.Context, cfg *config.Config, logger *slog.Logger) (http.Handler, error) {
		ps, err := pubsub.New(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create pubsub repository")
		}

		db, err := database.New(cfg, logger)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create database repository")
		}

		oa := &oauth.OAuth{
			GoogleConfig: &oauth2.Config{
				Endpoint:     google.Endpoint,
				ClientID:     cfg.UsersGoogleClientID,
				ClientSecret: cfg.UsersGoogleClientSecret,
				RedirectURL:  cfg.UsersRedirectURL,
				Scopes:       []string{"openid", "profile", "email"},
			},
		}

		svc := service.NewService(ctx, cfg, logger, ps, db, oa)
		return handler.NewHandler(cfg, logger, svc), nil
	}
}
