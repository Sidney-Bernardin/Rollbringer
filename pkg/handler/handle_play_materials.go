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

	var session, _ = r.Context().Value("session").(*domain.Session)

	// Create a PDF.
	pdf, err := h.Service.CreatePDF(r.Context(), &domain.PDF{
		OwnerID: session.UserID,
		GameID:  r.FormValue("game_id"),
		Name:    r.FormValue("name"),
		Schema:  r.FormValue("schema"),
	})

	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot create pdf"))
		return
	}

	// Respond with a PDFRow component.
	h.render(w, r, http.StatusOK, play_materials.PDFTableRow(pdf))
}

func (h *Handler) HandleGetPDFs(w http.ResponseWriter, r *http.Request) {

	var session, _ = r.Context().Value("session").(*domain.Session)

	// Get the PDFs.
	pdfs, err := h.Service.GetPDFs(r.Context(), session.UserID)
	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot get pdfs"))
		return
	}

	// Respond with a PDFs component.
	h.render(w, r, http.StatusOK, play_materials.PDFTableRows(pdfs))
}

func (h *Handler) HandleGetPDF(w http.ResponseWriter, r *http.Request) {

	// Get the PDF.
	pdf, err := h.Service.GetPDF(r.Context(), chi.URLParam(r, "pdf_id"))
	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot get pdf"))
		return
	}

	if chi.URLParam(r, "page_num") != "" {

		pageNum, err := strconv.Atoi(chi.URLParam(r, "page_num"))
		if err != nil {
			h.renderErr(w, r, http.StatusUnprocessableEntity, errors.New("page number must resemble a positive integer"))
			return
		}

		if pageNum-1 < 0 {
			h.renderErr(w, r, http.StatusUnprocessableEntity, errors.New("page number must resemble a positive integer"))
			return
		}

		h.render(w, r, http.StatusOK, components.PDFFields(pdf.ID, pdf.Pages[pageNum-1]))
		return
	}

	// Respond with a PDF-viewer tab.
	h.render(w, r, http.StatusOK, pages.PDFViewerTab(
		pdf.ID,
		pdf.Name,
		components.DNDCharacterSheetPageNames,
		components.DNDCharacterSheetFileLocation,
	))
}

func (h *Handler) HandleDeletePDF(w http.ResponseWriter, r *http.Request) {

	var session, _ = r.Context().Value("session").(*domain.Session)

	// Delete the PDF.
	if err := h.Service.DeletePDF(r.Context(), chi.URLParam(r, "pdf_id"), session.UserID); err != nil {
		h.err(w, r, errors.Wrap(err, "cannot delete pdf"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
