package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"rollbringer/internal"
	"rollbringer/internal/config"
)

var natsServer *server.Server

func StartEmbeddedServer(cfg *config.Config) (err error) {
	natsServer, err := server.NewServer(&server.Options{
		DontListen: !cfg.NATSEmbeddedServerListen,
		Host:       cfg.NATSHostname,
		Port:       cfg.NATSPort,
	})

	natsServer.ConfigureLogger()
	go natsServer.Start()
	if !natsServer.ReadyForConnections(cfg.NATSEmbeddedServerStartupTimeout) {
		return errors.New("server timed out on startup")
	}

	return nil
}

type PubSub struct {
	cfg    *config.Config
	logger *slog.Logger

	natsConn *nats.Conn
}

func New(cfg *config.Config, logger *slog.Logger) (*PubSub, error) {

	var (
		natsURL      = fmt.Sprintf("nats://%s:%v", cfg.NATSHostname, cfg.NATSPort)
		natsConnOpts = []nats.Option{}
	)

	if cfg.NATSEmbeddedServer && natsServer == nil {
		if err := StartEmbeddedServer(cfg); err != nil {
			return nil, errors.Wrap(err, "cannot start embedded NATS server")
		}
	}

	if natsServer != nil {
		natsURL = natsServer.ClientURL()
		natsConnOpts = append(natsConnOpts, nats.InProcessServer(natsServer))
	}
	fmt.Println(natsURL)

	natsConn, err := nats.Connect(natsURL, natsConnOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to NATS server")
	}

	return &PubSub{
		natsConn: natsConn,
		cfg:      cfg,
		logger:   logger,
	}, nil
}

func (ps *PubSub) Publish(ctx context.Context, subject string, event internal.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return internal.NewProblemDetail(ctx, &internal.PDOptions{
			Type:   internal.PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}

	err = ps.natsConn.Publish(subject, bytes)
	return errors.Wrap(err, "cannot publish event")
}

func (ps *PubSub) Request(ctx context.Context, subject string, req internal.Event, res any) error {
	reqBytes, err := internal.JSONEncodeEvent(ctx, req)
	if err != nil {
		return errors.Wrap(err, "cannot JSON encode event")
	}

	resMsg, err := ps.natsConn.RequestWithContext(ctx, subject, reqBytes)
	if err != nil {
		return errors.Wrap(err, "cannot request")
	}

	err = internal.JSONDecodeEvent(ctx, resMsg.Data, res)
	return errors.Wrap(err, "cannot JSON decode event")
}

func (ps *PubSub) Subscribe(
	ctx context.Context,
	subject string,
	errChan chan<- error,
	cb func(event internal.Event, subject []string) internal.Event,
) {
	var subCtx, cancel = context.WithCancel(ctx)
	defer cancel()

	sub, err := ps.natsConn.Subscribe(subject, func(msg *nats.Msg) {

		var req internal.Event
		if err := internal.JSONDecodeEvent(ctx, msg.Data, &req); err != nil {
			errChan <- errors.Wrap(err, "cannot JSON decode event")
		}

		res := cb(req, strings.Split(msg.Subject, "."))
		if res == nil {
			return
		}

		resBytes, err := internal.JSONEncodeEvent(ctx, res)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot JSON encode event")
		}

		if err := msg.Respond(resBytes); err != nil {
			errChan <- errors.Wrap(err, "cannot respond")
		}
	})

	if err != nil {
		errChan <- errors.Wrap(err, "cannot subscribe")
		return
	}

	sub.SetClosedHandler(func(subject string) {
		cancel()
	})
	defer sub.Unsubscribe()

	<-subCtx.Done()
}

func (ps *PubSub) ChanSubscribe(
	ctx context.Context,
	subject string,
	resChan chan<- internal.Event,
	errChan chan<- error,
) {
	var subCtx, cancel = context.WithCancel(ctx)
	defer cancel()

	sub, err := ps.natsConn.Subscribe(subject, func(msg *nats.Msg) {

		var event internal.Event
		if err := internal.JSONDecodeEvent(ctx, msg.Data, &event); err != nil {
			errChan <- errors.Wrap(err, "cannot JSON decode event")
		}

		resChan <- event
	})

	if err != nil {
		errChan <- errors.Wrap(err, "cannot subscribe")
		return
	}

	sub.SetClosedHandler(func(subject string) {
		cancel()
	})
	defer sub.Unsubscribe()

	<-subCtx.Done()
}
