package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/views/components"
)

func (api *API) handlePlayMaterials(conn *websocket.Conn) {

	var (
		r = conn.Request()

		incomingChan = make(chan domain.GameEvent)
		outgoingChan = make(chan domain.GameEvent)
	)

	go api.service.PlayMaterials(r.Context(), chi.URLParam(r, "game_id"), incomingChan, outgoingChan)

	go func() {
		defer api.closeConn(conn)

		for {
			select {
			case <-r.Context().Done():
				return

			case event := <-outgoingChan:
				if err := websocket.JSON.Send(conn, event); err != nil {
					api.err(conn, err, 0, wsStatusInternalError)
					return
				}
			}
		}
	}()

	for {

		var msg domain.GameEvent
		if err := websocket.JSON.Receive(conn, &msg); err != nil {

			if err == io.EOF {
				return
			}

			switch err.(type) {
			case *json.SyntaxError, *json.UnmarshalTypeError, *json.InvalidUnmarshalError:
				api.err(conn, err, 0, wsStatusUnsupportedData)
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
