package handler

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
	"rollbringer/pkg/views/pages"
)

func (h *Handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/play", http.StatusTemporaryRedirect)
}

func (h *Handler) HandlePlayPage(w http.ResponseWriter, r *http.Request) {

	var session, _ = r.Context().Value("session").(*domain.Session)

	// Get the play page.
	page, err := h.Service.GetPlayPage(r.Context(), session, r.URL.Query().Get("g"))
	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot get play page"))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "play_page", page))

	h.render(w, r, http.StatusOK, pages.Play())
}

func (h *Handler) HandleWebSocket(conn *websocket.Conn) {

	var (
		r = conn.Request()

		errChan      = make(chan error)
		incomingChan = make(chan domain.Event)
		outgoingChan = make(chan domain.Event)
	)

	// Process events in another go-routine.
	go h.Service.DoEvents(r.Context(), r.URL.Query().Get("g"), incomingChan, outgoingChan, errChan)

	// Respond with outoutgoing events in another go-routine.
	go func() {
		defer conn.Close()

		for {
			select {
			case <-r.Context().Done():
				return

			case err, ok := <-errChan:
				if !ok {
					return
				}

				h.renderErr(conn, r, 0, &domain.EventError{
					BaseEvent: domain.BaseEvent{
						Operation: "ERROR",
					},
					Err: err,
				})

			case e := <-outgoingChan:
				switch event := e.(type) {
				case *domain.EventUpdatePDFField:

					// Respond with the PDF field.
					h.render(conn, r, 0, components.PDFField(
						event.PDFID,
						event.FieldName,
						event.FieldValue,
					))

				default:
					h.Logger.Error().Any("event", e).Msg("Received event with unknown operation")
				}
			}
		}
	}()

	for {

		// Receive the next message from the client.
		var msg []byte
		if err := websocket.Message.Receive(conn, &msg); err != nil {
			if err == io.EOF || strings.Contains(err.Error(), net.ErrClosed.Error()) {
				return
			}

			errChan <- errors.Wrap(err, "cannot receive websocket message")
			return
		}

		var baseEvent domain.BaseEvent
		if err := json.Unmarshal(msg, &baseEvent); err != nil {
			errChan <- &domain.ProblemDetail{
				Type:   domain.PDTypeCannotDecodeEvent,
				Detail: err.Error(),
			}
			continue
		}

		event, err := baseEvent.GetOperationStruct()
		if err != nil {
			errChan <- errors.Wrap(err, "cannot get event operation struct")
			continue
		}

		if err := json.Unmarshal(msg, &event); err != nil {
			errChan <- &domain.ProblemDetail{
				Type:   domain.PDTypeCannotDecodeEvent,
				Detail: err.Error(),
			}
			continue
		}

		incomingChan <- event
	}
}
