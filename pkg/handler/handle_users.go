package handler

import (
	"net/http"
	"strings"
	"time"

	"rollbringer/pkg/domain"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {

	consentURL, state, codeVerifier := h.Service.StartLogin()

	// Store the state and code-verifier in a cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "STATE_AND_VERIFIER",
		Value:    state + "," + codeVerifier,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	http.Redirect(w, r, consentURL, http.StatusTemporaryRedirect)
}

func (h *Handler) HandleConsentCallback(w http.ResponseWriter, r *http.Request) {

	// Get the state/code-verifier cookie.
	cookie, err := r.Cookie("STATE_AND_VERIFIER")
	if err != nil {
		h.err(w, r, &domain.NormalError{
			Type: domain.NETypeUnauthorized,
		})
		return
	}

	// Get the state and code-verifier from the cookie.
	state_and_verifier := strings.Split(cookie.Value, ",")
	if len(state_and_verifier) != 2 {
		h.err(w, r, &domain.NormalError{
			Type: domain.NETypeUnauthorized,
		})
		return
	}

	session, err := h.Service.FinishLogin(r.Context(),
		state_and_verifier[0],
		r.FormValue("state"),
		r.FormValue("code"),
		state_and_verifier[1],
	)

	// Store the session-ID in a cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "SESSION_ID",
		Value:    session.ID.String(),
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
