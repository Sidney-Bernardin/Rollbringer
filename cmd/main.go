package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"rollbringer/src"
	"rollbringer/src/api"
	"rollbringer/src/domain/play"
	db_play "rollbringer/src/repositories/database/play"
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

	playDB, err := db_play.NewDatabase(config)
	if err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot create play-database", "err", err.Error())
		return
	}

	svr := api.NewServer(log, config,
		play.NewService(config, playDB, nil))

	/////

	go func() {
		log.Log(ctx, src.LevelInfo, "Listening", "address", config.APIAddr)
		err = svr.ListenAndServe()
		log.Log(ctx, src.LevelFatal, "Cannot listen and serve", "err", err.Error())
	}()

	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	<-signalCtx.Done()
	cancel()

	/////

	timeoutCtx, cancel := context.WithTimeout(ctx, config.GracfullShutdownTimeout)
	if err := svr.Shutdown(timeoutCtx); err != nil {
		log.Log(ctx, src.LevelFatal, "Cannot listen and serve", "err", err.Error())
	}
	cancel()
}
