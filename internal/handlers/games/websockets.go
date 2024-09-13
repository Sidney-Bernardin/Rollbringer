package games

import (
	"context"
	"encoding/json"
	"rollbringer/internal"
	"rollbringer/internal/views"
	"rollbringer/internal/views/pages/play"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

func (h *handler) handlePlayWebsocket(conn *websocket.Conn) {

	var (
		r   = conn.Request()
		ctx = r.Context()

		session, _ = r.Context().Value(internal.CtxKeySession).(*internal.Session)
		gameID, _  = uuid.Parse(chi.URLParam(r, "game_id"))

		pdfID             = uuid.Nil
		pdfCtx, pdfCancel = context.WithCancel(context.Background())

		resChan = make(chan any)
	)

	go func() {
		for {
			var bytes []byte
			if err := websocket.Message.Receive(conn, &bytes); err != nil {
				resChan <- errors.Wrap(err, "cannot receive WebSocket message")
				continue
			}

			var wrapper internal.EventWrapper[any]
			if err := json.Unmarshal(bytes, &wrapper); err != nil {
				resChan <- internal.NewProblemDetail(ctx, internal.PDOpts{
					Type:   internal.PDTypeInvalidJSON,
					Detail: err.Error(),
				})
				continue
			}

			var payload any
			switch wrapper.Event {
			case internal.EventSubToPDFRequest:
				payload = internal.SubToPDFRequest{}
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

			eventCtx := context.WithValue(ctx, internal.CtxKeyInstance, wrapper.Event)
			switch payload := payload.(type) {
			case internal.SubToPDFRequest:
				pdfCancel()

				pageFields, err := h.svc.GetPDFPage(eventCtx, payload.PDFID, payload.PageNum)
				if err != nil {
					resChan <- errors.Wrap(err, "cannot get PDF page")
					continue
				}

				pdfID = payload.PDFID
				pdfCtx, pdfCancel = context.WithCancel(eventCtx)

				go func() {
					err := h.svc.SubToPDFPage(pdfCtx, payload.PDFID, payload.PageNum, resChan)
					resChan <- errors.Wrap(err, "cannot sub to PDF page")
					pdfCancel()
				}()

				resChan <- &internal.PDFPage{
					PDFID:   pdfID,
					PageNum: payload.PageNum,
					Fields:  pageFields,
				}
			}
		}
	}()

	for {
		switch res := (<-resChan).(type) {
		case error:
			h.Err(conn, r, res)
		case *internal.PDFPage:
			h.Render(conn, r, 0, views.PDFViewerFields(res.PDFID, res.Fields))
		case *internal.Roll:
			h.Render(conn, r, 0, play.Roll(res))
		}
	}
}
