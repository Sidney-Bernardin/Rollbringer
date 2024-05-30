package handler

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views"
	"rollbringer/pkg/views/pages/play"
)

func (h *Handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/play", http.StatusTemporaryRedirect)
}

func (h *Handler) HandlePlayPage(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value("session").(*domain.Session)
		gameID, _  = uuid.Parse(r.URL.Query().Get("g"))
	)

	page, err := h.Service.GetPlayPage(r.Context(), session, gameID)
	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot get play page"))
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "play_page", page))

	h.render(w, r, http.StatusOK, play.Play())
}

func (h *Handler) HandleWebSocket(conn *websocket.Conn) {

	var (
		r              = conn.Request()
		ctx, cancelCtx = context.WithCancel(r.Context())

		incomingChan = make(chan domain.Event)
		outgoingChan = make(chan domain.Event)
		errChan      = make(chan error)

		gameID, _ = uuid.Parse(r.URL.Query().Get("g"))
	)

	defer cancelCtx()

	go h.Service.DoEvents(ctx, gameID, incomingChan, outgoingChan, errChan)

	go func() {
		defer cancelCtx()

		for {

			var eventBytes []byte
			if err := websocket.Message.Receive(conn, &eventBytes); err != nil {
				if err == io.EOF || strings.Contains(err.Error(), net.ErrClosed.Error()) {
					return
				}

				errChan <- errors.Wrap(err, "cannot receive websocket message")
				return
			}

			event, err := domain.DecodeJSONEvent(ctx, eventBytes)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot decode event")
				continue
			}

			incomingChan <- event
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case err, ok := <-errChan:
			if !ok {
				return
			}

			h.render(conn, r, 0, views.EventNormalError(&domain.EventNormalError{
				BaseEvent:   domain.BaseEvent{Operation: domain.OperationNormalError},
				NormalError: h.handleError(r.Context(), err),
			}))

		case e := <-outgoingChan:
			h.render(conn, r, 0, h.eventToComponent(e))
		}
	}
}

func (h *Handler) eventToComponent(e domain.Event) templ.Component {
	switch event := e.(type) {

	case *domain.EventPDFFields:
		return views.PDFViewerFields(event.PDFID, event.Fields)
	}

	return templ.NopComponent
}
