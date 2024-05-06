package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
)

// OpenIDConnectClaims represents an OpenID Connect JWT.
type OpenIDConnectClaims struct {
	*jwt.RegisteredClaims
	GivenName string `json:"given_name"`
}

func NewOauthState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type OAuth struct {
	GoogleConfig *oauth2.Config
}

func (oa *OAuth) NewCodeVerifier() string {
	return oauth2.GenerateVerifier()
}

func (oa *OAuth) GetConsentURL(state, codeVerifier string) string {
	return oa.GoogleConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(codeVerifier))
}

func (oa *OAuth) AuthenticateConsent(ctx context.Context, stateA, stateB, code, codeVerifier string) (*OpenIDConnectClaims, error) {

	// Verify the state.
	if stateA != stateB {
		return nil, &domain.ProblemDetail{
			Type: domain.PDTypeUnauthorized,
		}
	}

	// Exchange the code for an oauth-token.
	token, err := oa.GoogleConfig.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return nil, &domain.ProblemDetail{
			Type: domain.PDTypeUnauthorized,
		}
	}

	// Get the ID-token from the oauth-token.
	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token should be string, but is not")
	}

	// Parse the ID-token.
	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &OpenIDConnectClaims{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse ID token")
	}

	return idToken.Claims.(*OpenIDConnectClaims), nil
}
