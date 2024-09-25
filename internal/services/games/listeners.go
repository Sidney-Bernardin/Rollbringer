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

func (svc *service) Listen() error {
	errs, ctx := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		return svc.subToGames(ctx)
	})

	errs.Go(func() error {
		return svc.subToPDFs(ctx)
	})

	return errs.Wait()
}

func (svc *service) subToGames(ctx context.Context) error {
	err := svc.PubSub.Subscribe(ctx, "games", func(req *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch req.Event {
		case internal.EventGetGameRequest:
			payload = &internal.GetGameRequest{}
		case internal.EventGetGamesByHostRequest:
			payload = &internal.GetGamesByHostRequest{}
		case internal.EventGetGamesByUserRequest:
			payload = &internal.GetGamesByUserRequest{}
		default:
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidEvent,
					Detail: "The given evnet is not valid.",
					Extra: map[string]any{
						"event": req.Event,
					},
				}),
			}
		}

		instanceCtx := context.WithValue(ctx, internal.CtxKeyInstance, req.Event)

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(instanceCtx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				}),
			}
		}

		switch payload := payload.(type) {
		case *internal.GetGameRequest:
			game, err := svc.getGame(instanceCtx, payload.GameID, payload.ViewQuery)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(instanceCtx, svc.Logger, errors.Wrap(err, "cannot get game")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventGame,
				Payload: game,
			}

		case *internal.GetGamesByHostRequest:
			games, err := svc.getGamesByHost(instanceCtx, payload.HostID, payload.ViewQuery)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(instanceCtx, svc.Logger, errors.Wrap(err, "cannot get games by host")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventGames,
				Payload: games,
			}

		case *internal.GetGamesByUserRequest:
			games, err := svc.getGamesByUser(instanceCtx, payload.UserID, payload.ViewQuery)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(instanceCtx, svc.Logger, errors.Wrap(err, "cannot get games by guest")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventGames,
				Payload: games,
			}

		default:
			return &internal.EventWrapper[any]{
				Event:   internal.EventError,
				Payload: internal.HandleError(ctx, svc.Logger, internal.SvrErrUnknownEvent),
			}
		}
	})

	return errors.Wrap(err, "cannot subscribe to games")
}

func (svc *service) subToPDFs(ctx context.Context) error {
	err := svc.PubSub.Subscribe(ctx, "pdfs", func(req *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch req.Event {
		case internal.EventGetPDFsByOwnerRequest:
			payload = &internal.GetPDFsByOwnerRequest{}
		default:
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidEvent,
					Detail: "The given evnet is not valid.",
					Extra: map[string]any{
						"event": req.Event,
					},
				}),
			}
		}

		instanceCtx := context.WithValue(ctx, internal.CtxKeyInstance, req.Event)

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return &internal.EventWrapper[any]{
				Event: internal.EventError,
				Payload: internal.NewProblemDetail(instanceCtx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				}),
			}
		}

		switch payload := payload.(type) {
		case *internal.GetPDFsByOwnerRequest:
			pdfs, err := svc.getPDFsByOwner(instanceCtx, payload.OwnerID, payload.ViewQuery)
			if err != nil {
				return &internal.EventWrapper[any]{
					Event:   internal.EventError,
					Payload: internal.HandleError(instanceCtx, svc.Logger, errors.Wrap(err, "cannot get PDFs by owner")),
				}
			}

			return &internal.EventWrapper[any]{
				Event:   internal.EventPDFs,
				Payload: pdfs,
			}

		default:
			return &internal.EventWrapper[any]{
				Event:   internal.EventError,
				Payload: internal.HandleError(ctx, svc.Logger, internal.SvrErrUnknownEvent),
			}
		}
	})

	return errors.Wrap(err, "cannot subscribe to games")
}

func (svc *service) SubToGame(ctx context.Context, gameID uuid.UUID, resChan chan<- any) error {
	subject := fmt.Sprintf("games.%s", gameID)
	err := svc.PubSub.Subscribe(ctx, subject, func(req *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch req.Event {
		case internal.EventPDFPage:
			payload = &internal.PDFPage{}
		case internal.EventRoll:
			payload = &internal.Roll{}
		default:
			return nil
		}

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
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

	return errors.Wrapf(err, "cannot subscribe to games.%s", gameID)
}

func (svc *service) SubToPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, resChan chan<- any) error {
	subject := fmt.Sprintf("pdfs.%s.pages.%v", pdfID, pageNum)
	err := svc.PubSub.Subscribe(ctx, subject, func(req *internal.EventWrapper[[]byte]) *internal.EventWrapper[any] {

		var payload any
		switch req.Event {
		case internal.EventPDFPage:
			payload = &internal.PDFPage{}
		default:
			return nil
		}

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
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

	return errors.Wrapf(err, "cannot subscribe to pdfs.%s.pages.%v", pdfID, pageNum)
}
