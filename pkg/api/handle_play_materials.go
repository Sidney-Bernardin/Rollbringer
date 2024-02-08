package api

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
	"rollbringer/pkg/views/components/play_materials"
	"rollbringer/pkg/views/oob_swaps"
)

func (api *API) handlePlayMaterials(conn *websocket.Conn) {

	var (
		r = conn.Request()

		incomingChan = make(chan *domain.GameEvent)
		outgoingChan = make(chan *domain.GameEvent)
	)

	go api.service.PlayMaterials(r.Context(), r.URL.Query().Get("g"), incomingChan, outgoingChan)

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
					api.render(conn, r, oob_swaps.UpdatePDFFields(event), 0)
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

		var event *domain.GameEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			api.err(conn, err, 0, wsStatusUnsupportedData)
			return
		}

		incomingChan <- event
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
