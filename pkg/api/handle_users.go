package api

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
)

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {

	var (
		state        = mustGetRandHexStr()
		codeVerifier = oauth2.GenerateVerifier()
	)

	// Store the state and code-verifier in a cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "STATE_AND_VERIFIER",
		Value:    state + "," + codeVerifier,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	// Generate and redirect to the consent URL.
	consentURL := api.GoogleOAuthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(codeVerifier))
	http.Redirect(w, r, consentURL, http.StatusTemporaryRedirect)
}

func (api *API) HandleConsentCallback(w http.ResponseWriter, r *http.Request) {

	// Get the state/code-verifier cookie.
	cookie, err := r.Cookie("STATE_AND_VERIFIER")
	if err != nil {
		api.err(w, domain.ErrUnauthorized, http.StatusUnauthorized, 0)
		return
	}

	// Get the state and code-verifier from the cookie.
	state_and_verifier := strings.Split(cookie.Value, ",")
	if len(state_and_verifier) != 2 {
		api.err(w, domain.ErrUnauthorized, http.StatusUnauthorized, 0)
		return
	}

	// Verify both state.
	if r.FormValue("state") != state_and_verifier[0] {
		api.err(w, domain.ErrUnauthorized, http.StatusUnauthorized, 0)
		return
	}

	// Exchange the code for an oauth token.
	token, err := api.GoogleOAuthConfig.Exchange(
		r.Context(),
		r.FormValue("code"),
		oauth2.VerifierOption(state_and_verifier[1]))

	if err != nil {
		api.err(w, domain.ErrUnauthorized, http.StatusUnauthorized, 0)
		return
	}

	// Get the ID-token from the oauth token.
	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		err = errors.New("id_token should be string, but is not")
		api.err(w, err, http.StatusInternalServerError, 0)
		return
	}

	// Parse the ID-token.
	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &openIDConnectClaims{})
	if err != nil {
		err = errors.Wrap(err, "cannot parse ID token")
		api.err(w, err, http.StatusInternalServerError, 0)
		return
	}

	// Login the user.
	sessionID, err := api.Service.LoginUser(r.Context(), idToken.Claims.(*openIDConnectClaims).Subject)
	if err != nil {
		api.domainErr(w, errors.Wrap(err, "cannot login user"))
		return
	}

	// Store the session-ID in a cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "SESSION_ID",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
