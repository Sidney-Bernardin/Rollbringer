package api

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func (api *API) HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {

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

func (api *API) HandleOAuthConsentCallback(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("state_and_verifier")
	if err != nil {
		// errUnauthorized
		return
	}

	state_and_verifier := strings.Split(cookie.Value, ",")
	if len(state_and_verifier) != 2 {
		// errUnauthorized
		return
	}

	if r.FormValue("state") != state_and_verifier[0] {
		// errUnauthorized
		return
	}

	token, err := api.GoogleOAuthConfig.Exchange(
		r.Context(),
		r.FormValue("code"),
		oauth2.VerifierOption(state_and_verifier[1]))

	if err != nil {
		err = errors.Wrap(err, "cannot exchange code for token")
		// server error
		return
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		err = errors.New("id_token should be string, but is not")
		// server error
		return
	}

	_, _, err = jwt.NewParser().ParseUnverified(idTokenStr, &openIDConnectClaims{})
	if err != nil {
		err = errors.Wrap(err, "cannot parse ID token")
		// server error
		return
	}

	// login
	// set session token cookie

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
