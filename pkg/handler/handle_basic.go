package handler

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/oob_swaps"
	"rollbringer/pkg/views/pages"
)

func (h *Handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/play", http.StatusTemporaryRedirect)
}

func (h *Handler) HandlePlayPage(w http.ResponseWriter, r *http.Request) {

	// Get the game.
	game, err := h.Service.GetGame(r.Context(), r.URL.Query().Get("g"))
	if err != nil && errors.Cause(err) != domain.ErrGameNotFound {
		h.domainErr(w, errors.Wrap(err, "cannot get game"))
		return
	}
	giveToRequest(r, "game", game)

	// Check if the user is logged in by getting the session. If the user is
	// logged out, render the page early.
	session, _ := r.Context().Value("session").(*domain.Session)
	if session == nil {
		h.render(w, r, pages.Play(), http.StatusOK)
		return
	}

	// Get the user.
	user, err := h.Service.GetUser(r.Context(), session.UserID)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get user"))
		return
	}
	giveToRequest(r, "user", user)

	// Get the user's games.
	games, err := h.Service.GetGamesFromUser(r.Context(), session.UserID)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get user's games"))
		return
	}
	giveToRequest(r, "games", games)

	// Get the user's PDFs.
	pdfs, err := h.Service.GetPDFs(r.Context(), session.UserID)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get pdfs"))
		return
	}
	giveToRequest(r, "pdfs", pdfs)

	h.render(w, r, pages.Play(), http.StatusOK)
}

func (h *Handler) HandleWebSocket(conn *websocket.Conn) {

	var (
		r = conn.Request()

		incomingChan = make(chan *domain.Event)
		outgoingChan = make(chan *domain.Event)
	)

	go h.Service.DoEvents(r.Context(), r.URL.Query().Get("g"), incomingChan, outgoingChan)

	go func() {
		defer conn.Close()

		for {
			select {
			case <-r.Context().Done():
				return

			case event, ok := <-outgoingChan:

				if !ok {
					return
				}

				switch event.Type {
				case "UPDATE_PDF_FIELDS":
					h.render(conn, r, oob_swaps.UpdatePDFFields(event), 0)
				}
			}
		}
	}()

	for {
		var msg []byte
		if err := websocket.Message.Receive(conn, &msg); err != nil {
			if err == io.EOF || strings.Contains(err.Error(), net.ErrClosed.Error()) {
				return
			}

			h.err(conn, err, 0, wsStatusInternalError)
			return
		}

		var event domain.Event
		if err := json.Unmarshal(msg, &event); err != nil {
			h.err(conn, err, 0, wsStatusUnsupportedData)
			return
		}

		incomingChan <- &event
	}
}
