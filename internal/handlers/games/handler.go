package games

import (
	"fmt"
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
	"rollbringer/internal/views/games"
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
	h.Router.With(h.Authenticate("", false)).Method("GET", "/ws", websocket.Handler(h.handleGameWebsocket))

	h.Router.Route("/", func(r chi.Router) {
		r.Use(h.Authenticate("", true))

		r.Post("/games", h.handleCreateGame)
		r.Delete("/games/{game_id}", h.handleDeleteGame)

		r.Post("/pdfs", h.handleCreatePDF)
		r.Get("/pdfs/{pdf_id}", h.handleGetPDF)
		r.Put("/pdfs/{pdf_id}", h.handleUpdatePDF)
		r.Delete("/pdfs/{pdf_id}", h.handleDeletePDF)
	})

	return h
}

func (h *handler) handleCreateGame(w http.ResponseWriter, r *http.Request) {

	var (
		ctx = r.Context()

		session, _ = ctx.Value(internal.CtxKeySession).(*internal.Session)
		game       = &internal.Game{
			Name: r.FormValue("name"),
		}
	)

	if err := h.svc.CreateGame(ctx, session, game); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot create game"))
		return
	}

	h.Render(w, r, http.StatusCreated, games.HostedGameRow(game))
}

func (h *handler) handleDeleteGame(w http.ResponseWriter, r *http.Request) {

	var (
		ctx = r.Context()

		session, _ = ctx.Value(internal.CtxKeySession).(*internal.Session)
		gameID, _  = uuid.Parse(chi.URLParam(r, "game_id"))
	)

	if err := h.svc.DeleteGame(ctx, session, gameID); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot delete game"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) handleCreatePDF(w http.ResponseWriter, r *http.Request) {

	var (
		ctx        = r.Context()
		session, _ = ctx.Value(internal.CtxKeySession).(*internal.Session)
	)

	gameID, err := internal.OptionalID(ctx, r.FormValue("game_id"))
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot parse game-id"))
		return
	}

	pdf := &internal.PDF{
		GameID: gameID,
		Name:   r.FormValue("name"),
		Schema: r.FormValue("schema"),
	}

	if err = h.svc.CreatePDF(ctx, session, pdf); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot create pdf"))
		return
	}

	pdf, err = h.svc.GetPDF(ctx, pdf.ID, internal.PDFView(r.URL.Query().Get("view")))
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot get pdf"))
		return
	}

	h.Render(w, r, http.StatusCreated, games.PDFTableRow(pdf, true))
}

func (h *handler) handleGetPDF(w http.ResponseWriter, r *http.Request) {
	var pdfID, _ = uuid.Parse(chi.URLParam(r, "pdf_id"))

	pdf, err := h.svc.GetPDF(r.Context(), pdfID, "")
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot get pdf"))
		return
	}

	h.Render(w, r, http.StatusOK, games.PDFTab(pdf))
}

func (h *handler) handleUpdatePDF(w http.ResponseWriter, r *http.Request) {

	var (
		ctx = r.Context()

		session, _ = ctx.Value(internal.CtxKeySession).(*internal.Session)
		pdfID, _   = uuid.Parse(chi.URLParam(r, "pdf_id"))
		view       = internal.PDFView(r.URL.Query().Get("view"))
	)

	err := h.svc.UpdatePDF(ctx, session, &internal.PDF{
		ID:   pdfID,
		Name: r.FormValue("name"),
	})

	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot update pdf"))
		return
	}

	pdf, err := h.svc.GetPDF(ctx, pdfID, view)
	if err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot get pdf"))
		return
	}

	h.Render(w, r, http.StatusOK, games.PDFTableRow(pdf, true))
}

func (h *handler) handleDeletePDF(w http.ResponseWriter, r *http.Request) {

	var (
		ctx = r.Context()

		session, _ = ctx.Value(internal.CtxKeySession).(*internal.Session)
		pdfID, _   = uuid.Parse(chi.URLParam(r, "pdf_id"))
	)

	if err := h.svc.DeletePDF(ctx, session, pdfID); err != nil {
		h.Err(w, r, errors.Wrap(err, "cannot delete pdf"))
		return
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf(`remove-tab-%s, deleted-pdf-%s`, pdfID, pdfID))
	w.WriteHeader(http.StatusOK)
}
