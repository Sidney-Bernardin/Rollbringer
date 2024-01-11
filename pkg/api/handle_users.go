package api

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {

	var (
		state        = mustGetRandHexStr()
		codeVerifier = oauth2.GenerateVerifier()
	)

	http.SetCookie(w, &http.Cookie{
		Name:     "state_and_verifier",
		Value:    state + "," + codeVerifier,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	consentURL := api.GoogleOAuthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(codeVerifier))
	http.Redirect(w, r, consentURL, http.StatusTemporaryRedirect)
}

func (api *API) HandleConsentCallback(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("state_and_verifier")
	if err != nil {
		api.err(w, r, errUnauthorized, http.StatusUnauthorized)
		return
	}

	state_and_verifier := strings.Split(cookie.Value, ",")
	if len(state_and_verifier) != 2 {
		api.err(w, r, errUnauthorized, http.StatusUnauthorized)
		return
	}

	if r.FormValue("state") != state_and_verifier[0] {
		api.err(w, r, errUnauthorized, http.StatusUnauthorized)
		return
	}

	token, err := api.GoogleOAuthConfig.Exchange(
		r.Context(),
		r.FormValue("code"),
		oauth2.VerifierOption(state_and_verifier[1]))

	if err != nil {
		err = errors.Wrap(err, "cannot exchange code for token")
		api.err(w, r, err, http.StatusUnauthorized)
		return
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		err = errors.New("id_token should be string, but is not")
		api.err(w, r, err, http.StatusInternalServerError)
		return
	}

	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &openIDConnectClaims{})
	if err != nil {
		err = errors.Wrap(err, "cannot parse ID token")
		api.err(w, r, err, http.StatusInternalServerError)
		return
	}

	session, err := api.DB.Login(r.Context(), idToken.Claims.(*openIDConnectClaims).Subject)
	if err != nil {
		api.err(w, r, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "SESSION_ID",
		Value:    session.ID.String(),
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
