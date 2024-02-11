package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
	"rollbringer/pkg/views/components/play_materials"
	"rollbringer/pkg/views/oob_swaps"
)

func (h *Handler) HandleCreatePDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Create a PDF.
	pdfID, name, err := h.Service.CreatePDF(r.Context(), session.UserID, r.FormValue("schame"))
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot create pdf"))
		return
	}

	// Respond with a PDFRow component.
	h.render(w, r, play_materials.PDFRow(pdfID, name, "foo bar"), http.StatusOK)
}

func (h *Handler) HandleGetPDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Get the PDF.
	pdf, err := h.Service.GetPDF(r.Context(), chi.URLParam(r, "pdf_id"), session.UserID)
	if err != nil {
		h.domainErr(w, errors.Wrap(err, "cannot get pdf"))
		return
	}

	// Respond with a PDF-viewer tab.
	h.render(w, r,
		oob_swaps.AddPDFViewerTab(
			pdf.Name,
			components.PDFViewer(
				pdf.ID,
				components.DNDCharacterSheetFileLocation,
				components.DNDCharacterSheetPageNames,
			),
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
