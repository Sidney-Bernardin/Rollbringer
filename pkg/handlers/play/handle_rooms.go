package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/handlers/play/views"
)

func (h *playHandler) handleRoomsPost(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		room = &domain.Room{
			Name: r.FormValue("name"),
		}
	)

	if err := h.playSvc.CreateRoom(ctx, state["session"].(*domain.Session), room); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot create room", nil))
		return
	}

	h.Respond(w, r, http.StatusOK, views.RoomListItem(room))
}

func (h *playHandler) handleRoomsDelete(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		roomID, _ = uuid.Parse(chi.URLParam(r, "room_id"))
	)

	if err := h.playSvc.DeleteRoom(ctx, state["session"].(*domain.Session), roomID); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot delete room", nil))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *playHandler) handleRoomsWebSocket(conn *websocket.Conn) {

	var (
		r = conn.Request()
		// state = h.State(r)
		// ctx   = r.Context()

		// session   = state["session"].(*domain.Session)
		// roomID, _ = uuid.Parse(chi.URLParam(r, "room_id"))

		resChan = make(chan *domain.Event)
	)

	go func() {
		for {

			// Receive WebSocket message.
			var msg []byte
			if err := websocket.Message.Receive(conn, &msg); err != nil {
				h.Err(conn, r, domain.Wrap(err, "cannot receive message", nil))
				return
			}
			fmt.Println(string(msg))

			// Decode the message.
			event, err := domain.DecodeEvent(msg, map[domain.Operation]any{})
			if err != nil {
				h.Err(conn, r, domain.Wrap(err, "cannot decode event", nil))
				return
			}

			switch event.Payload.(type) {
			}
		}
	}()

	for {
		switch ((<-resChan).Payload).(type) {
		}
	}
}
