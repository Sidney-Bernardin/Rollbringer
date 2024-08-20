package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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

func (ps *PubSub) Publish(ctx context.Context, subject string, data internal.Event) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return &internal.ProblemDetail{
			Instance: ctx.Value(internal.CtxKeyInstance).(string),
			Type:     internal.PDTypeInvalidEvent,
			Detail:   err.Error(),
		}
	}

	err = ps.natsConn.Publish(subject, bytes)
	return errors.Wrap(err, "cannot publish event")
}

func (ps *PubSub) Request(ctx context.Context, subject string, data []byte) ([]byte, error) {
	msg, err := ps.natsConn.RequestWithContext(ctx, subject, data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot request")
	}
	return msg.Data, nil
}

func (ps *PubSub) Subscribe(
	ctx context.Context,
	subject string,
	errChan chan<- error,
	cb func(msg internal.Event) internal.Event,
) {
	var subCtx, cancel = context.WithCancel(ctx)
	defer cancel()

	sub, err := ps.natsConn.Subscribe(subject, func(msg *nats.Msg) {
		event, err := internal.JSONDecodeEvent(ctx, msg.Data)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot JSON decode event")
		}

		res := cb(event)
		if res == nil {
			return
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			errChan <- &internal.ProblemDetail{
				Instance: ctx.Value(internal.CtxKeyInstance).(string),
				Type:     internal.PDTypeInvalidEvent,
				Detail:   err.Error(),
			}
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
		event, err := internal.JSONDecodeEvent(ctx, msg.Data)
		if err != nil {
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
