package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"rollbringer/src"
	"rollbringer/src/api"
	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/domain/services/play"
	"rollbringer/src/repos/broker"
	accounts_database "rollbringer/src/repos/database/accounts"
	play_database "rollbringer/src/repos/database/play"
	"rollbringer/src/repos/google"
	"rollbringer/src/repos/spotify"
)

var (
	logJSON = flag.Bool("log_json", false, "")
	noColor = flag.Bool("no_color", false, "")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	var log *slog.Logger
	if *logJSON {
		log = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	} else {
		log = src.NewPrettyLogger(*noColor)
	}

	config, err := src.NewConfig()
	if err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot create configuration", "err", err.Error())
		return
	}

	publicBroker, err := broker.New(ctx, config, log)
	if err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot create public-broker", "err", err.Error())
		return
	}

	//
	///// Accounts

	accountsDatabase, err := accounts_database.NewDatabase(ctx, config)
	if err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot create accounts-database", "err", err.Error())
		return
	}

	google := google.New(config)
	spotify := spotify.New(config)

	accountsSvc := accounts.NewService(config, publicBroker, accountsDatabase, google, spotify)

	//
	///// Play

	playBroker, err := broker.NewPlayBroker(ctx, publicBroker.(*broker.PublicBroker))
	if err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot create play-broker", "err", err.Error())
		return
	}

	playDatabase, err := play_database.NewDatabase(ctx, config)
	if err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot create play-database", "err", err.Error())
		return
	}

	playSvc := play.NewService(config, log, playBroker, playDatabase)

	//
	///// Server

	svr := api.NewServer(log, config, publicBroker,
		accountsSvc, accountsDatabase, google, spotify,
		playSvc, playBroker, playDatabase)

	go func() {
		log.Log(ctx, src.LevelInfo, "Listening", "address", config.APIAddr)
		err = svr.ListenAndServe()
		log.Log(ctx, src.LevelFatal, "Cannot listen and serve", "err", err.Error())
	}()

	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	<-signalCtx.Done()
	cancel()

	//
	///// Shutdown

	timeoutCtx, cancel := context.WithTimeout(ctx, config.GracfullShutdownTimeout)
	if err := svr.Shutdown(timeoutCtx); err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot listen and serve", "err", err.Error())
	}
	cancel()
}
