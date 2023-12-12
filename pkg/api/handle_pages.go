package api

import (
	"net/http"
	"rollbringer/pkg/database"
)

type defaultPageTmpl struct {
	Name      string `json:"name"`
	CSRFToken string `json:"-"`
}

type errorPageTmpl struct {
	defaultPageTmpl

	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`

	err error
}

func newErrorPageTmpl(err error) (ret errorPageTmpl) {
	ret.Name = "Error"
	ret.Message = err.Error()
	ret.err = err
	return ret
}

func (a *api) handleHomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageData := defaultPageTmpl{Name: "Home"}

		if session, ok := r.Context().Value("session").(*database.Session); ok {
			pageData.CSRFToken = session.CSRFToken
		}

		a.executeTemplate(w, "page.html", http.StatusOK, pageData)
	}
}

func (a *api) handleGamePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageData := defaultPageTmpl{Name: "DND"}

		if session, ok := r.Context().Value("session").(*database.Session); ok {
			pageData.CSRFToken = session.CSRFToken
		}

		a.executeTemplate(w, "page.html", http.StatusOK, pageData)
	}
}
