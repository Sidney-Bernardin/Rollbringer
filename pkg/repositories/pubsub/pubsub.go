package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

var embeddedServer *server.Server

func createEmbeddedServer(config *domain.Config, logger *slog.Logger) (err error) {

	// Crate NATS server.
	embeddedServer, err = server.NewServer(&server.Options{
		DontListen: !config.NATSEmbeddedServerListen,
		Host:       config.NATSHostname,
		Port:       config.NATSPort,
	})

	if err != nil {
		return domain.Wrap(err, "cannot create NATS server", nil)
	}

	// Set the logger for the server.
	embeddedServer.SetLoggerV2(&natsLogger{logger}, false, false, false)

	// Start server.
	go embeddedServer.Start()
	if !embeddedServer.ReadyForConnections(5 * time.Second) {
		return errors.New("timed out on startup")
	}

	return nil
}

type PubSubRepository struct {
	config *domain.Config
	logger *slog.Logger

	conn *nats.Conn
}

func NewPubSubRepository(config *domain.Config, logger *slog.Logger) (domain.PubSubRepository, error) {

	var (
		url  = fmt.Sprintf("nats://%s:%v", config.NATSHostname, config.NATSPort)
		opts = []nats.Option{}
	)

	// Create an embedded NATS server if it's configured and doesn't already exist.
	if config.NATSEmbeddedServer && embeddedServer == nil {
		if err := createEmbeddedServer(config, logger); err != nil {
			return nil, domain.Wrap(err, "cannot create embedded NATS server", nil)
		}
	}

	// If the embedded server exists, use it.
	if embeddedServer != nil {
		url = embeddedServer.ClientURL()
		opts = append(opts, nats.InProcessServer(embeddedServer))
	}

	// Create connection to NATS server.
	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, domain.Wrap(err, "cannot connect to NATS server", nil)
	}

	return &PubSubRepository{config, logger, conn}, nil
}

func (repo *PubSubRepository) Close() {
	repo.conn.Close()
}

func (repo *PubSubRepository) Publish(ctx context.Context, subject string, event *domain.Event) error {

	// Encode the event's payload.
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return domain.Wrap(err, "cannot JSON encode request event payload", nil)
	}

	// Create a message for the event.
	msg := nats.NewMsg(subject)
	msg.Header.Add("operation", string(event.Operation))
	msg.Data = payload

	// Publish the message.
	err = repo.conn.PublishMsg(msg)
	return domain.Wrap(err, "cannot JSON encode request event payload", nil)
}

func (repo *PubSubRepository) Request(ctx context.Context, subject string, resPayload any, reqEvent *domain.Event) (domain.Operation, error) {

	// Encode the request-event's payload.
	reqPayload, err := json.Marshal(reqEvent.Payload)
	if err != nil {
		return "", domain.Wrap(err, "cannot JSON encode request event payload", nil)
	}

	// Create a message for the request-event.
	reqMsg := nats.NewMsg(subject)
	reqMsg.Header.Add("operation", string(reqEvent.Operation))
	reqMsg.Data = reqPayload

	// Send request-message.
	resMsg, err := repo.conn.RequestMsgWithContext(ctx, reqMsg)
	if err != nil {
		return "", domain.Wrap(err, "cannot send request", map[string]any{
			"subject": subject,
		})
	}

	resOperation := domain.Operation(resMsg.Header.Get("operation"))

	// Treat the response-payload as a UserError if the response-operation is ERROR.
	var userErr *domain.UserError
	if resOperation == domain.OperationError {
		resPayload = userErr
	}

	// Decode the response-message's data.
	if err := json.Unmarshal(resMsg.Data, resPayload); err != nil {
		return "", domain.Wrap(err, "cannot JSON decode response event payload", nil)
	}

	return resOperation, userErr
}

func (repo *PubSubRepository) Subscribe(ctx context.Context, subject string, dpFunc domain.DefinePayloadFunc, handlerFunc func(*domain.Event) *domain.Event) error {

	// Subscribe to the subject.
	subChan := make(chan *nats.Msg, 1)
	sub, err := repo.conn.ChanSubscribe(subject, subChan)
	if err != nil {
		return domain.Wrap(err, "cannot subscribe", map[string]any{
			"subject": subject,
		})
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil

		case reqMsg := <-subChan:
			go func() {

				var (
					reqOperation = domain.Operation(reqMsg.Header.Get("operation"))
					reqPayload   = dpFunc(reqOperation)
				)

				// Decode the request-message's data.
				if err := json.Unmarshal(reqMsg.Data, &reqPayload); err != nil {
					domain.HandleError(ctx, repo.logger, slog.LevelError, domain.Wrap(err, "cannot JSON decode event payload", nil))
					return
				}

				// Call the handler function.
				resEvent := handlerFunc(&domain.Event{
					Operation: reqOperation,
					Payload:   reqPayload,
				})

				if resEvent == nil {
					return
				}

				// Encode the response-event's payload.
				resPayload, err := json.Marshal(resEvent.Payload)
				if err != nil {
					domain.HandleError(ctx, repo.logger, slog.LevelError, domain.Wrap(err, "cannot JSON encode event payload", nil))
					return
				}

				// Create a message for the response-event.
				resMsg := nats.NewMsg(reqMsg.Reply)
				resMsg.Header.Add("operation", string(resEvent.Operation))
				resMsg.Data = resPayload

				// Send response-message.
				if err := reqMsg.RespondMsg(resMsg); err != nil {
					domain.HandleError(ctx, repo.logger, slog.LevelError, domain.Wrap(err, "cannot send response", map[string]any{
						"subject": subject,
					}))
				}
			}()
		}
	}
}
