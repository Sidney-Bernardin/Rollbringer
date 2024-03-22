package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
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
	go h.Service.DoEvents(r.Context(), r.URL.Query().Get("g"), errChan, incomingChan, outgoingChan)

	// Prepare outoutgoing events in another go-routine.
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

				h.renderErr(conn, r, 0, domain.NewEventError(err))

			case e := <-outgoingChan:
				switch event := e.(type) {
				case *domain.EventUpdatePDFField:

					// Render the PDF field.
					h.render(conn, r, 0, components.PDFField(
						event.PDFID,
						event.Headers.HXTrigger,
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

			h.renderErr(conn, r, 0, domain.NewEventError(errors.Wrap(err, "cannot receive websocket message")))
			return
		}

		var e domain.BaseEvent
		if err := json.Unmarshal(msg, &e); err != nil {
			h.renderErr(conn, r, 0, domain.NewEventError(&domain.ProblemDetail{
				Type:   domain.PDTypeCannotDecodeEvent,
				Detail: err.Error(),
			}))
			continue
		}

		event, ok := domain.OperationTypes[e.Operation]
		if !ok {
			h.renderErr(conn, r, 0, domain.NewEventError(&domain.ProblemDetail{
				Type:   domain.PDTypeInvalidEventOperation,
				Detail: fmt.Sprintf(`"%s" is an invlid event operation`, e.Operation),
			}))
			continue
		}
		event = reflect.New(reflect.TypeOf(event)).Interface().(domain.Event)

		if err := json.Unmarshal(msg, &event); err != nil {
			fmt.Println(err)
			h.renderErr(conn, r, 0, domain.NewEventError(&domain.ProblemDetail{
				Type:   domain.PDTypeCannotDecodeEvent,
				Detail: err.Error(),
			}))
			continue
		}

		incomingChan <- event
	}
}
