package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
	"rollbringer/pkg/views/components/play_materials"
	"rollbringer/pkg/views/oob_swaps"
)

func (api *API) handleCreatePDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Create a PDF.
	pdfID, name, err := api.service.CreatePDF(r.Context(), session.UserID, r.FormValue("schame"))
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot create pdf"))
		return
	}

	// Respond with a PDFRow component.
	api.render(w, r, play_materials.PDFRow(pdfID, name, "foo bar"), http.StatusOK)
}

func (api *API) handleGetPDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Get the PDF.
	pdf, err := api.service.GetPDF(r.Context(), chi.URLParam(r, "pdf_id"), session.UserID)
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot get pdf"))
		return
	}

	// Respond with a PDF-viewer tab.
	api.render(w, r,
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

func (api *API) handleDeletePDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Delete the PDF.
	if err := api.service.DeletePDF(r.Context(), chi.URLParam(r, "pdf_id"), session.UserID); err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot delete pdf"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
