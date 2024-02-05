package api

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
)

func (api *API) handlePlayMaterials(conn *websocket.Conn) {

	var (
		r = conn.Request()

		incomingChan = make(chan []byte)
		outgoingChan = make(chan *domain.GameEvent)
	)

	go api.service.PlayMaterials(r.Context(), r.URL.Query().Get("g"), incomingChan, outgoingChan)

	go func() {
		defer conn.Close()

		for {
			select {
			case <-r.Context().Done():
				return

			case event := <-outgoingChan:

				if event == nil {
					return
				}

				switch event.Type {
				case "UPDATE_PDF_PAGE":
					api.render(conn, r, components.ReplacePDFFields(event), 0)
				case "INIT_PDF_PAGE":
					api.render(conn, r, components.ReplacePDFFields(event), 0)
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

			api.err(conn, err, 0, wsStatusInternalError)
			return
		}

		incomingChan <- msg
	}
}

func (api *API) handleCreatePDF(w http.ResponseWriter, r *http.Request) {

	session, _ := r.Context().Value("session").(*domain.Session)

	// Create a PDF.
	pdfID, name, err := api.service.CreatePDF(r.Context(), session.UserID, r.FormValue("schame"))
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot create pdf"))
		return
	}

	// Respond with a PDFRow component.
	api.render(w, r, components.PDFRow(pdfID, name, "foo bar"), http.StatusOK)
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
	api.render(w, r, components.AddPDFViewerTab(pdf.Name, components.DNDCharacterSheet(pdf.ID)), http.StatusOK)
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
