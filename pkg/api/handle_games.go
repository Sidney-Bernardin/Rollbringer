package api

import (
	"bytes"

	"golang.org/x/net/websocket"
)

func (a *api) handleWS() websocket.Handler {

	type updatedPDFField struct {
		TextArea bool
		Type     string
		Name     string
		Value    string
	}

	return func(conn *websocket.Conn) {
		for {
			var msg string
			if err := websocket.Message.Receive(conn, &msg); err != nil {
				a.logger.Error().Stack().Err(err).Msg("Cannot recive ws msg")
				return
			}

			x := updatedPDFField{
				TextArea: false,
				Type:     "text",
				Name:     "testfield",
				Value:    "new value",
			}

			bbuf := bytes.Buffer{}
			if err := a.tmpl.ExecuteTemplate(&bbuf, "updated_pdf_fields", x); err != nil {
				a.logger.Error().Stack().Err(err).Msg("Cannot send ws template")
				return
			}

			conn.Write(bbuf.Bytes())
		}
	}
}
