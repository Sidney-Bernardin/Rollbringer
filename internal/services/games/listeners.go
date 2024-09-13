package games

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"rollbringer/internal"
)

func (svc *service) Run() error {
	errs, ctx := errgroup.WithContext(context.Background())
	go svc.subToGames(ctx)
	return errs.Wait()
}

func (svc *service) subToGames(ctx context.Context) error {
	err := svc.PS.Subscribe(ctx, "games", func(incoming *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch incoming.Event {
		case internal.EventGetGameRequest:
			payload = internal.GetGameRequest{}
		default:
			return nil
		}

		if err := json.Unmarshal(incoming.Payload, &payload); err != nil {
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				}),
			}
		}

		switch payload := payload.(type) {
		case internal.GetGameRequest:
			game, err := svc.getGame(ctx, payload.GameID, payload.View)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(ctx, svc.Logger, errors.Wrap(err, "cannot get game")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventGame,
				Payload: game,
			}

		default:
			return nil
		}
	})

	return errors.Wrap(err, "cannot subscribe to games")
}

func (svc *service) SubToPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, resChan chan<- any) error {
	subject := fmt.Sprintf("pdfs.%s.pages.%v", pdfID, pageNum)
	err := svc.PS.Subscribe(ctx, subject, func(incoming *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch incoming.Event {
		case internal.EventPDFPage:
			payload = internal.PDFPage{}
		default:
			return nil
		}

		if err := json.Unmarshal(incoming.Payload, &payload); err != nil {
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				}),
			}
		}

		resChan <- payload
		return nil
	})

	return errors.Wrap(err, "cannot subscribe to pdfs.%s.pages.%v")
}
