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
		resChan           = make(chan any)
	)

	gameID, err := internal.OptionalID(ctx, r.URL.Query().Get("g"))
	if err != nil {
		h.Err(conn, r, errors.Wrap(err, "cannot parse game-ID"))
		return
	}

	if gameID != nil {
		go func() {
			err := h.svc.SubToGame(ctx, *gameID, resChan)
			resChan <- errors.Wrap(err, "cannot subscribe to game")
		}()
	}

	go func() {
		for {
			var bytes []byte
			if err := websocket.Message.Receive(conn, &bytes); err != nil {
				if err == io.EOF {
					return
				}

				resChan <- errors.Wrap(err, "cannot receive WebSocket message")
				continue
			}

			var req internal.EventWrapper[any]
			if err := json.Unmarshal(bytes, &req); err != nil {
				resChan <- internal.NewProblemDetail(ctx, internal.PDOpts{
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
			default:
				continue
			}

			if err := json.Unmarshal(bytes, &payload); err != nil {
				resChan <- internal.NewProblemDetail(ctx, internal.PDOpts{
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
					resChan <- errors.Wrap(err, "cannot get PDF page")
					continue
				}

				pdfCtx, pdfCancel = context.WithCancel(instanceCtx)

				go func() {
					err := h.svc.SubToPDFPage(pdfCtx, payload.PDFID, payload.PageNum, resChan)
					if errors.Cause(err) != context.Canceled {
						resChan <- errors.Wrap(err, "cannot subscribe to PDF page")
					}
					pdfCancel()
				}()

				resChan <- &internal.PDFPage{
					PDFID:   payload.PDFID,
					PageNum: payload.PageNum,
					Fields:  pageFields,
				}

			case *internal.UpdatePDFPageRequest:
				if err := h.svc.UpdatePDFPage(ctx, payload.PDFID, payload.PageNum, payload.FieldName, payload.FieldValue); err != nil {
					resChan <- errors.Wrap(err, "cannot update PDF page")
				}

			case *internal.CreateRollRequest:
				if gameID == nil {
					continue
				}

				if err := h.svc.CreateRoll(ctx, session, *gameID, payload.DiceTypes, payload.Modifiers); err != nil {
					resChan <- errors.Wrap(err, "cannot create roll")
				}
			}
		}
	}()

	for {
		switch res := (<-resChan).(type) {
		case error:
			h.Err(conn, r, res)
		case *internal.PDFPage:
			h.Render(conn, r, 0, games.PDFViewerFields(res.PDFID, res.Fields))
		case *internal.Roll:
			h.Render(conn, r, 0, games.Roll(res))
		}
	}
}
