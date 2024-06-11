package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/pages/play"
)

func (h *Handler) HandleCreatePDF(w http.ResponseWriter, r *http.Request) {

	var session, _ = r.Context().Value("session").(*domain.Session)

	view, ok := domain.PDFViews[r.URL.Query().Get("view")]
	if !ok {
		h.err(w, r, &domain.NormalError{
			Type:   domain.NETypeInvalidView,
			Detail: "The given PDF view is unknown.",
		})
	}

	pdf := &domain.PDF{
		OwnerID: session.UserID,
		Name:    r.FormValue("name"),
		Schema:  r.FormValue("schema"),
	}

	if gameID, err := uuid.Parse(r.FormValue("game_id")); err == nil {
		pdf.GameID = &gameID
	}

	err := h.Service.CreatePDF(r.Context(), session, pdf, view)
	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot create pdf"))
		return
	}

	h.render(w, r, http.StatusOK, play.NewPDFTableRow(pdf))
}

func (h *Handler) HandleGetPDF(w http.ResponseWriter, r *http.Request) {

	var pdfID, _ = uuid.Parse(chi.URLParam(r, "pdf_id"))

	pdf, err := h.Service.GetPDF(r.Context(), pdfID, domain.PDFViewAll)
	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot get pdf"))
		return
	}

	h.render(w, r, http.StatusOK, play.PDFTab(pdf))
}

func (h *Handler) HandleDeletePDF(w http.ResponseWriter, r *http.Request) {

	var (
		session, _ = r.Context().Value("session").(*domain.Session)
		pdfID, _   = uuid.Parse(chi.URLParam(r, "pdf_id"))
	)

	if err := h.Service.DeletePDF(r.Context(), session, pdfID); err != nil {
		h.err(w, r, errors.Wrap(err, "cannot delete pdf"))
		return
	}

	h.render(w, r, http.StatusOK, play.DeletePDF(pdfID))
}
