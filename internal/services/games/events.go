package games

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (svc *service) ProcessGameEvents(
	ctx context.Context,
	session *internal.Session,
	gameID uuid.UUID,
	incomingChan <-chan []byte,
	outgoingChan chan<- []byte,
	errChan chan<- error,
) {
	defer close(errChan)

	var (
		pdfID        uuid.UUID
		pdfCtx       context.Context
		pdfCtxCancel context.CancelFunc = func() {}
	)

	if gameID != uuid.Nil {
		_, err := svc.db.GameGet(ctx, gameID, internal.GameViewAll)
		if err != nil {
			errChan <- errors.Wrap(err, "cannot get game")
			return
		}

		go svc.subscribeToGame(ctx, gameID, incomingChan, errChan)
	}

	for {
		select {
		case <-ctx.Done():
			pdfCtxCancel()
			return

		case e := <-incomingChan:
			if err := e.Validate(ctx); err != nil {
				errChan <- errors.Wrap(err, "invalid event")
				continue
			}

			switch event := e.(type) {
			case *internal.EventSubToPDF:
				pdfCtxCancel()
				var eventCtx = context.WithValue(ctx, internal.CtxKeyInstance, internal.ETSubToPDF)

				pdfFields, err := svc.db.PDFGetPage(eventCtx, event.PDFID, event.PageNum-1)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot get PDF fields")
					continue
				}

				pdfID = event.PDFID
				pdfCtx, pdfCtxCancel = context.WithCancel(eventCtx)

				go svc.PS.ChanSubscribe(pdfCtx, fmt.Sprintf("pdfs.%s", pdfID), outgoingChan, errChan)

				outgoingChan <- &internal.EventPDFFields{
					BaseEvent: internal.BaseEvent{Type: internal.ETPdfFields},
					PDFID:     pdfID,
					PageNum:   event.PageNum,
					Fields:    pdfFields,
				}

			case *internal.EventUpdatePDFField:
				var eventCtx = context.WithValue(ctx, internal.CtxKeyInstance, internal.ETUpdatePDFField)

				if pdfID == uuid.Nil {
					errChan <- internal.NewProblemDetail(eventCtx, internal.PDOpts{
						Type:   internal.PDTypeNotSubscribedToPDF,
						Detail: "You must be subscribed to a PDF before updating it.",
					})
				}

				err := svc.db.PDFUpdatePage(eventCtx, pdfID, event.PageNum-1, event.FieldName, event.FieldValue)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot update pdf page")
					continue
				}

				err = svc.PS.Publish(eventCtx, "pdfs."+pdfID.String(), &internal.EventPDFFields{
					BaseEvent: internal.BaseEvent{Type: internal.ETPdfFields},
					PDFID:     pdfID,
					PageNum:   event.PageNum,
					Fields: map[string]string{
						event.FieldName: event.FieldValue,
					},
				})

				if err != nil {
					errChan <- errors.Wrap(err, "cannot publish event")
				}

			case *internal.EventCreateRoll:
				var eventCtx = context.WithValue(ctx, internal.CtxKeyInstance, internal.ETCreateRoll)

				roll, err := svc.roll(eventCtx, event.Dice)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot roll")
					continue
				}
				roll.OwnerID = session.UserID
				roll.GameID = gameID

				if err := svc.db.RollInsert(eventCtx, roll); err != nil {
					errChan <- errors.Wrap(err, "cannot insert roll")
					continue
				}

				err = svc.PS.Publish(eventCtx, "games."+gameID.String(), &internal.EventRoll{
					BaseEvent: internal.BaseEvent{Type: internal.ETRoll},
					Roll:      *roll,
				})

				if err != nil {
					errChan <- errors.Wrap(err, "cannot publish event")
				}
			}
		}
	}
}
