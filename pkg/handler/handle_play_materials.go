package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
	"rollbringer/pkg/views/components/play_materials"
	"rollbringer/pkg/views/pages"
)

func (h *Handler) HandleCreatePDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Create a PDF.
	pdf, err := h.Service.CreatePDF(r.Context(), &domain.PDF{
		OwnerID: session.UserID,
		GameID:  r.FormValue("game_id"),
		Name:    r.FormValue("name"),
		Schema:  r.FormValue("schema"),
	})

	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot create pdf"))
		return
	}

	// Respond with a PDFRow component.
	h.render(w, r, play_materials.PDFTableRow(pdf), http.StatusOK)
}

func (h *Handler) HandleGetPDFs(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Get the PDFs.
	pdfs, err := h.Service.GetPDFs(r.Context(), session.UserID)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get pdfs"))
		return
	}

	// Respond with a PDFs component.
	h.render(w, r, play_materials.PDFTableRows(pdfs), http.StatusOK)
}

func (h *Handler) HandleGetPDF(w http.ResponseWriter, r *http.Request) {

	_, _ = r.Context().Value("session").(*domain.Session)

	// Get the PDF.
	pdf, err := h.Service.GetPDF(r.Context(), chi.URLParam(r, "pdf_id"))
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get pdf"))
		return
	}

	if chi.URLParam(r, "page_num") != "" {

		pageNum, err := strconv.Atoi(chi.URLParam(r, "page_num"))
		if err != nil {
			h.err(w, errors.New("page number must resemble a positive integer"), http.StatusUnprocessableEntity, 0)
			return
		}

		if pageNum-1 < 0 {
			h.err(w, errors.New("page number must resemble a positive integer"), http.StatusUnprocessableEntity, 0)
			return
		}

		h.render(w, r, components.PDFFields(pdf.ID, pdf.Pages[pageNum-1]), http.StatusOK)
		return
	}

	// Respond with a PDF-viewer tab.
	h.render(w, r,
		pages.PDFViewerTab(
			pdf.ID,
			pdf.Name,
			components.DNDCharacterSheetPageNames,
			components.DNDCharacterSheetFileLocation,
		),
		http.StatusOK,
	)
}

func (h *Handler) HandleDeletePDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Delete the PDF.
	if err := h.Service.DeletePDF(r.Context(), chi.URLParam(r, "pdf_id"), session.UserID); err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot delete pdf"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
