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

func (ps *PubSub) Publish(ctx context.Context, subject string, data *internal.EventWrapper[any]) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeInvalidJSON,
			Detail: err.Error(),
		})
	}

	err = ps.natsConn.Publish(subject, bytes)
	return errors.Wrap(err, "cannot publish")
}

func (ps *PubSub) Request(ctx context.Context, subject string, res any, req *internal.EventWrapper[any]) error {
	reqPayload, err := json.Marshal(req.Payload)
	if err != nil {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeInvalidJSON,
			Detail: err.Error(),
		})
	}

	reqMsg := nats.NewMsg(subject)
	reqMsg.Header.Add("event", string(req.Event))
	reqMsg.Data = reqPayload

	resMsg, err := ps.natsConn.RequestMsgWithContext(ctx, reqMsg)
	if err != nil {
		return errors.Wrap(err, "cannot request")
	}

	resEvent := internal.Event(resMsg.Header.Get("event"))

	if resEvent == internal.EventError {
		var pd internal.ProblemDetail
		if err := json.Unmarshal(resMsg.Data, &pd); err != nil {
			return errors.Wrap(err, "cannot JSON decode problem-detail")
		}
		return &pd
	}

	if res, ok := res.(*internal.EventWrapper[[]byte]); ok {
		res.Event = resEvent
		res.Payload = resMsg.Data
		return nil
	}

	err = json.Unmarshal(resMsg.Data, res)
	return errors.Wrap(err, "cannot JSON decode response")
}

func (ps *PubSub) Subscribe(ctx context.Context, subject string, cb func(*internal.EventWrapper[[]byte]) *internal.EventWrapper[any]) error {
	subscriptionCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sub, err := ps.natsConn.Subscribe(subject, func(reqMsg *nats.Msg) {
		res := cb(&internal.EventWrapper[[]byte]{
			Event:   internal.Event(reqMsg.Header.Get("event")),
			Payload: reqMsg.Data,
		})

		if res == nil {
			return
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			internal.HandleError(subscriptionCtx, ps.logger, errors.Wrap(err, "cannot JSON encode response"))
			return
		}

		resMsg := nats.NewMsg("")
		resMsg.Header.Add("event", string(res.Event))
		resMsg.Data = resBytes

		if err := reqMsg.RespondMsg(resMsg); err != nil {
			internal.HandleError(subscriptionCtx, ps.logger, errors.Wrap(err, "cannot respond"))
		}
	})

	if err != nil {
		return errors.Wrap(err, "cannot subscribe")
	}

	defer sub.Unsubscribe()
	sub.SetClosedHandler(func(subject string) {
		cancel()
	})

	<-subscriptionCtx.Done()
	return subscriptionCtx.Err()
}
