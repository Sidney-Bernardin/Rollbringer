package games

import (
	"context"
	"encoding/json"
	"io"
	"rollbringer/internal"
	"rollbringer/internal/views/games"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

func (h *handler) handleGameWebsocket(conn *websocket.Conn) {

	var (
		r   = conn.Request()
		ctx = r.Context()

		session, _ = ctx.Value(internal.CtxKeySession).(*internal.Session)

		pdfCtx, pdfCancel = context.WithCancel(context.Background())

		resChan = make(chan *internal.EventWrapper[any])
		errChan = make(chan error)
	)

	gameID, err := internal.OptionalID(ctx, r.URL.Query().Get("g"))
	if err != nil {
		pdfCancel()
		h.Err(conn, r, errors.Wrap(err, "cannot parse game-ID"))
		return
	}

	if gameID != nil {
		go func() {
			err := h.svc.SubToGame(ctx, *gameID, resChan)
			errChan <- errors.Wrap(err, "cannot subscribe to game")
		}()
	}

	go func() {
		for {
			var bytes []byte
			if err := websocket.Message.Receive(conn, &bytes); err != nil {
				if err == io.EOF {
					return
				}

				errChan <- errors.Wrap(err, "cannot receive WebSocket message")
				continue
			}

			var req internal.EventWrapper[any]
			if err := json.Unmarshal(bytes, &req); err != nil {
				errChan <- internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				})
				continue
			}

			var payload any
			switch req.Event {
			case internal.EventSubToPDFRequest:
				payload = &internal.SubToPDFRequest{}
			case internal.EventUpdatePDFPageRequest:
				payload = &internal.UpdatePDFPageRequest{}
			case internal.EventCreateRollRequest:
				payload = &internal.CreateRollRequest{}
			case internal.EventCreateChatMessageRequest:
				payload = &internal.CreateChatMessageRequest{}
			default:
				continue
			}

			if err := json.Unmarshal(bytes, &payload); err != nil {
				errChan <- internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				})
				continue
			}

			instanceCtx := context.WithValue(ctx, internal.CtxKeyInstance, req.Event)
			switch payload := payload.(type) {
			case *internal.SubToPDFRequest:
				pdfCancel()

				pageFields, err := h.svc.GetPDFPage(instanceCtx, payload.PDFID, payload.PageNum)
				if err != nil {
					errChan <- errors.Wrap(err, "cannot get PDF page")
					continue
				}

				pdfCtx, pdfCancel = context.WithCancel(instanceCtx)

				go func() {
					err := h.svc.SubToPDFPage(pdfCtx, payload.PDFID, payload.PageNum, resChan)
					if errors.Cause(err) != context.Canceled {
						errChan <- errors.Wrap(err, "cannot subscribe to PDF page")
					}
					pdfCancel()
				}()

				resChan <- &internal.EventWrapper[any]{
					Event: internal.EventPDFPage,
					Payload: &internal.PDFPage{
						PDFID:   payload.PDFID,
						PageNum: payload.PageNum,
						Fields:  pageFields,
					},
				}

			case *internal.UpdatePDFPageRequest:
				if err := h.svc.UpdatePDFPage(ctx, payload.PDFID, payload.PageNum, payload.FieldName, payload.FieldValue, true); err != nil {
					errChan <- errors.Wrap(err, "cannot update PDF page")
				}

			case *internal.CreateRollRequest:
				if gameID == nil {
					continue
				}

				if _, err := h.svc.CreateRoll(ctx, session, *gameID, payload.DiceTypes, payload.Modifiers); err != nil {
					errChan <- errors.Wrap(err, "cannot create roll")
					continue
				}

			case *internal.CreateChatMessageRequest:
				if gameID == nil {
					continue
				}

				if _, err := h.svc.CreateChatMessage(ctx, session, *gameID, payload.Message); err != nil {
					errChan <- errors.Wrap(err, "cannot create chat-message")
					continue
				}
			}
		}
	}()

	for {
		select {
		case err := <-errChan:
			h.Err(conn, r, err)

		case event := <-resChan:
			switch payload := event.Payload.(type) {
			case *internal.PDF:

				switch event.Event {
				case internal.EventDeletedPDF:
					h.Render(conn, r, 0, event)
				default:
					h.Render(conn, r, 0, games.PDFTableRowOOB(payload, payload.OwnerID == session.UserID))
				}

			case *internal.PDFPage:
				h.Render(conn, r, 0, games.PDFViewerFields(payload.PDFID, payload.Fields))
			case *internal.Roll:
				h.Render(conn, r, 0, games.Roll(payload))
			case *internal.ChatMessage:
				h.Render(conn, r, 0, games.ChatMessage(payload))
			}
		}
	}
}
