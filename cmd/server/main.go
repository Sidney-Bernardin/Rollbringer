package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	HTTP "github.com/Sidney-Bernardin/Rollbringer/server/http"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/google"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/nats"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/sql"
	"github.com/Sidney-Bernardin/Rollbringer/server/service"
	"golang.org/x/oauth2"
	oaGoogle "golang.org/x/oauth2/google"
)

var (
	jsonLogs = flag.Bool("json", false, "Log JSON.")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	var log *slog.Logger
	if *jsonLogs {
		log = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	} else {
		log = server.NewPrettyLogger()
	}

	config, err := server.NewConfig()
	if err != nil {
		log.Log(ctx, slog.LevelError, "Cannot load configuration", "err", err.Error())
		return
	}

	/////

	sql, err := sql.New(ctx, config, log.With("namespace", "sql"))
	if err != nil {
		log.Log(ctx, slog.LevelError, "Cannot create SQL repository", "err", err.Error())
		return
	}

	nats, err := nats.New(ctx, config)
	if err != nil {
		log.Log(ctx, slog.LevelError, "Cannot create repositories", "err", err.Error())
		return
	}

	google := &google.Google{
		OAuthConfig: &oauth2.Config{
			Endpoint:     oaGoogle.Endpoint,
			ClientID:     config.GoogleOauthClientId,
			ClientSecret: config.GoogleOauthClientSecret,
			RedirectURL:  config.GoogleOauthRedirectUrl,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}

	/////

	api := &HTTP.API{
		Server: &http.Server{
			Addr: config.HttpAddr,
		},
		Log: log.With("namespace", "http"),
		Service: &service.Service{
			Config: config,
			Log:    log.With("namespace", "service"),
			SQL:    sql,
			Nats:   nats,
			Google: google,
		},
	}

	go func() {
		api.Server.Handler = api.Router()
		if err := api.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Log(ctx, slog.LevelError, "HTTP failed", "err", err.Error())
		}
	}()

	/////

	sigCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	<-sigCtx.Done()
	cancel()

	timeoutCtx, cancel := context.WithTimeout(ctx, config.HttpGracfullShutdownTimeout)
	if err := api.Shutdown(timeoutCtx); err != nil {
		log.Log(ctx, slog.LevelError, "Gracefull shutdown failed", "err", err.Error())
	}
	cancel()
}
