package games

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/handlers"
	service "rollbringer/internal/services/games"
	"rollbringer/internal/views/pages/play"
)

type handler struct {
	*handlers.BaseHandler

	svc service.Service
}

func NewHandler(cfg *config.Config, logger *slog.Logger, svc service.Service) http.Handler {
	h := &handler{
		BaseHandler: &handlers.BaseHandler{
			Config:      cfg,
			Logger:      logger,
			Router:      chi.NewRouter(),
			BaseService: svc,
		},
		svc: svc,
	}

	h.Router.Use(h.Log, h.Instance)
	h.Router.With(h.GetSession).Method("GET", "/ws", websocket.Handler(h.handleGameWebsocket))

	authRouter := h.Router.With(h.Authenticate)

	authRouter.Post("/games", h.handleCreateGame)
	authRouter.Delete("/games/{game_id}", h.handleDeleteGame)

	authRouter.Post("/pdfs", h.handleCreatePDF)
	authRouter.Get("/pdfs/{pdf_id}", h.handleGetPDF)
	authRouter.Delete("/pdfs/{pdf_id}", h.handleDeletePDF)

	return h
}

func (h *handler) handleCreateGame(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value(internal.CtxKeySession).(*internal.Session)
		game       = &internal.Game{
			Name: r.FormValue("name"),
		}
	)

	if err := h.svc.CreateGame(r.Context(), session, game); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot create game"))
		return
	}

	h.Render(w, r, http.StatusCreated, play.HostedGameRow(game))
}

func (h *handler) handleDeleteGame(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value(internal.CtxKeySession).(*internal.Session)
		gameID, _  = uuid.Parse(chi.URLParam(r, "game_id"))
	)

	if err := h.svc.DeleteGame(r.Context(), session, gameID); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot delete game"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) handleCreatePDF(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value(internal.CtxKeySession).(*internal.Session)
		view       = r.URL.Query().Get("view")
		pdf        = &internal.PDF{
			Name:   r.FormValue("name"),
			Schema: r.FormValue("schema"),
		}
	)

	if gameID, err := uuid.Parse(r.FormValue("game_id")); err == nil {
		pdf.GameID = &gameID
	}

	err := h.svc.CreatePDF(r.Context(), session, pdf, "pdf-all,owner-all,game-all")
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot create pdf"))
		return
	}

	h.Render(w, r, http.StatusCreated, play.NewPDFTableRow(pdf, view == "with-game-row"))
}

func (h *handler) handleGetPDF(w http.ResponseWriter, r *http.Request) {
	var pdfID, _ = uuid.Parse(chi.URLParam(r, "pdf_id"))

	pdf, err := h.svc.GetPDF(r.Context(), pdfID, "pdf-all")
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot get pdf"))
		return
	}

	h.Render(w, r, http.StatusOK, play.PDFTab(pdf))
}

func (h *handler) handleDeletePDF(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value(internal.CtxKeySession).(*internal.Session)
		pdfID, _   = uuid.Parse(chi.URLParam(r, "pdf_id"))
	)

	if err := h.svc.DeletePDF(r.Context(), session, pdfID); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot delete pdf"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
