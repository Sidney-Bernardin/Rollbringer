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
		return svc.subToPages(ctx)
	})

	return errs.Wait()
}

func (svc *service) subToPages(ctx context.Context) error {
	err := svc.PubSub.Subscribe(ctx, "games.pages", func(req *internal.EventWrapper[[]byte]) (*internal.EventWrapper[any], error) {

		var payload any
		switch req.Event {
		case internal.EventPlayPage:
			payload = &internal.PlayPage{}
		default:
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeInvalidEvent,
				Detail: "The given evnet is not valid.",
				Extra: map[string]any{
					"event": req.Event,
				},
			})
		}

		instanceCtx := context.WithValue(ctx, internal.CtxKeyInstance, req.Event)

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, internal.NewProblemDetail(instanceCtx, internal.PDOpts{
				Type:   internal.PDTypeInvalidJSON,
				Detail: err.Error(),
			})
		}

		switch payload := payload.(type) {
		case *internal.PlayPage:
			err := svc.playPage(instanceCtx, payload)
			return &internal.EventWrapper[any]{
				Event:   internal.EventGames,
				Payload: payload,
			}, errors.Wrap(err, "cannot get page data")

		default:
			return nil, internal.HandleError(ctx, svc.Logger, errors.New("invalid games.pages event"))
		}
	})

	return errors.Wrap(err, "cannot subscribe to games.pages")
}

func (svc *service) SubToGame(ctx context.Context, gameID uuid.UUID, resChan chan<- *internal.EventWrapper[any]) error {
	subject := fmt.Sprintf("games.%s", gameID)
	err := svc.PubSub.Subscribe(ctx, subject, func(req *internal.EventWrapper[[]byte]) (*internal.EventWrapper[any], error) {

		var payload any
		switch req.Event {
		case internal.EventPDF:
			payload = &internal.PDF{}
		case internal.EventDeletedPDF:
			payload = &internal.PDF{}
		case internal.EventPDFPage:
			payload = &internal.PDFPage{}
		case internal.EventRoll:
			payload = &internal.Roll{}
		default:
			return nil, nil
		}

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeInvalidJSON,
				Detail: err.Error(),
			})
		}

		resChan <- &internal.EventWrapper[any]{
			Event:   req.Event,
			Payload: payload,
		}

		return nil, nil
	})

	return errors.Wrapf(err, "cannot subscribe to games.%s", gameID)
}

func (svc *service) SubToPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, resChan chan<- *internal.EventWrapper[any]) error {
	subject := fmt.Sprintf("pdfs.%s.pages.%v", pdfID, pageNum)
	err := svc.PubSub.Subscribe(ctx, subject, func(req *internal.EventWrapper[[]byte]) (*internal.EventWrapper[any], error) {

		var payload any
		switch req.Event {
		case internal.EventPDFPage:
			payload = &internal.PDFPage{}
		default:
			return nil, nil
		}

		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type:   internal.PDTypeInvalidJSON,
				Detail: err.Error(),
			})
		}

		resChan <- &internal.EventWrapper[any]{
			Event:   req.Event,
			Payload: payload,
		}

		return nil, nil
	})

	return errors.Wrapf(err, "cannot subscribe to pdfs.%s.pages.%v", pdfID, pageNum)
}
