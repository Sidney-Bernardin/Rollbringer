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

func newEmbeddedServer(cfg *config.Config) (*server.Server, error) {
	svr, err := server.NewServer(&server.Options{
		DontListen: !cfg.NATSEmbeddedServerListen,
		Host:       cfg.NATSHostname,
		Port:       cfg.NATSPort,
	})

	if err != nil {
		return nil, errors.Wrap(err, "cannot create NATS server")
	}

	svr.ConfigureLogger()
	go svr.Start()
	if !svr.ReadyForConnections(cfg.NATSEmbeddedServerStartupTimeout) {
		return nil, errors.New("server timed out on startup")
	}

	return svr, nil
}

type PubSub struct {
	cfg    *config.Config
	logger *slog.Logger

	natsConn *nats.Conn
}

func New(cfg *config.Config, logger *slog.Logger) (internal.PubSub, error) {

	var (
		natsURL      = fmt.Sprintf("nats://%s:%v", cfg.NATSHostname, cfg.NATSPort)
		natsConnOpts = []nats.Option{}
	)

	if cfg.NATSEmbeddedServer && natsServer == nil {
		var err error
		natsServer, err = newEmbeddedServer(cfg)
		if err != nil {
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
		cfg:      cfg,
		logger:   logger,
		natsConn: natsConn,
	}, nil
}

func (ps *PubSub) Close() {
	ps.natsConn.Close()
}

func (ps *PubSub) Publish(ctx context.Context, subject string, event internal.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeInvalidEvent,
			Detail: err.Error(),
		})
	}

	err = ps.natsConn.Publish(subject, bytes)
	return errors.Wrap(err, "cannot publish event")
}

func (ps *PubSub) Request(ctx context.Context, subject string, req internal.Event) (internal.Event, error) {
	reqBytes, err := internal.JSONEncodeEvent(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot JSON encode event")
	}

	resMsg, err := ps.natsConn.RequestWithContext(ctx, subject, reqBytes)
	if err != nil {
		return nil, errors.Wrap(err, "cannot request")
	}

	res, err := internal.JSONDecodeEvent(ctx, resMsg.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot JSON decode event")
	}

	if event, ok := res.(*internal.EventError); ok {
		return nil, &event.ProblemDetail
	}

	return res, nil
}

func (ps *PubSub) Subscribe(
	ctx context.Context,
	subject string,
	errChan chan<- error,
	cb func(event internal.Event, subject []string) (internal.Event, *internal.ProblemDetail),
) {
	sub, err := ps.natsConn.Subscribe(subject, func(msg *nats.Msg) {
		req, err := internal.JSONDecodeEvent(ctx, msg.Data)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot JSON decode event")
		}

		res, pd := cb(req, strings.Split(msg.Subject, "."))
		if pd != nil {
			resBytes, err := internal.JSONEncodeEvent(ctx, &internal.EventError{
				BaseEvent:     internal.BaseEvent{Type: internal.ETError},
				ProblemDetail: *pd,
			})

			if err != nil {
				errChan <- errors.Wrap(err, "cannot JSON encode event")
			}

			if err := msg.Respond(resBytes); err != nil {
				errChan <- errors.Wrap(err, "cannot respond")
			}
		} else if res != nil {
			resBytes, err := internal.JSONEncodeEvent(ctx, res)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot JSON encode event")
			}

			if err := msg.Respond(resBytes); err != nil {
				errChan <- errors.Wrap(err, "cannot respond")
			}
		}
	})

	if err != nil {
		errChan <- errors.Wrap(err, "cannot subscribe")
		return
	}

	defer sub.Unsubscribe()
	<-ctx.Done()
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
