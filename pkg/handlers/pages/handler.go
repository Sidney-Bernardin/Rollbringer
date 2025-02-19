package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/handlers"
	"rollbringer/pkg/handlers/pages/views"
)

type pagesHandler struct {
	*handlers.Handler

	pubSub domain.PubSubRepository
}

func New(config *domain.Config, logger *slog.Logger, svc *domain.Service) http.Handler {
	h := &pagesHandler{
		Handler: &handlers.Handler{
			Config:  config,
			Logger:  logger,
			Router:  chi.NewRouter(),
			Service: svc,
		},
		pubSub: svc.PubSub,
	}

	h.Router.Use(h.MWLog)
	h.Router.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(os.DirFS("cmd/static"))))
	h.Router.With(h.MWAuthenticate(false, false, true)).Get("/", h.handleHomePage)
	h.Router.With(h.MWAuthenticate(true, false, true)).Get("/play", h.handlePlayPage)

	return h
}

func (h *pagesHandler) handleHomePage(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		session, _ = state["session"].(*domain.Session)
	)

	if session != nil {
		_, err := h.pubSub.Request(ctx, "play", &session.User.Rooms, &domain.Event{
			Operation: domain.OperationGetRoomsRequest,
			Payload: domain.GetRoomsRequest{
				OwnerID: session.UserID,
			},
		})

		if err != nil {
			h.Err(w, r, domain.Wrap(err, "cannot get rooms", nil))
			return
		}
	}

	h.Respond(w, r, http.StatusOK, views.HomePage(&views.HomePageState{
		Session: session,
	}))
}

func (h *pagesHandler) handlePlayPage(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		session   = state["session"].(*domain.Session)
		roomID, _ = uuid.Parse(r.URL.Query().Get("r"))
	)

	var room *domain.Room
	_, err := h.pubSub.Request(ctx, "play", &room, &domain.Event{
		Operation: domain.OperationGetRoomRequest,
		Payload: domain.GetRoomRequest{
			RoomID: roomID,
		},
	})

	if err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot get rooms", nil))
		return
	}

	h.Respond(w, r, http.StatusOK, views.PlayPage(&views.PlayPageState{
		Session: session,
		Room:    room,
	}))
}
