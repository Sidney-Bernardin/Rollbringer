package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"rollbringer/internal"
	"rollbringer/internal/config"
)

// openIDConnectClaims represents an OpenID Connect JWT.
type openIDConnectClaims struct {
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

func New(cfg *config.Config) *OAuth {
	return &OAuth{
		GoogleConfig: &oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     cfg.UsersGoogleClientID,
			ClientSecret: cfg.UsersGoogleClientSecret,
			RedirectURL:  cfg.UsersRedirectURL,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}
}

func (oa *OAuth) GenerateCodeVerifier() string {
	return oauth2.GenerateVerifier()
}

func (oa *OAuth) GetConsentURL(state, codeVerifier string) string {
	return oa.GoogleConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(codeVerifier))
}

func (oa *OAuth) AuthenticateConsent(ctx context.Context, stateA, stateB, code, codeVerifier string) (*internal.GoogleUserInfo, error) {

	// Verify the state.
	if stateA != stateB {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeUnauthorized,
		})
	}

	// Exchange the code for an oauth-token.
	token, err := oa.GoogleConfig.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeUnauthorized,
		})
	}

	// Get the ID-token from the oauth-token.
	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token should be string, but is not")
	}

	// Parse the ID-token.
	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &openIDConnectClaims{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse ID token")
	}

	claims, _ := idToken.Claims.(*openIDConnectClaims)
	return &internal.GoogleUserInfo{
		GoogleID:  claims.Subject,
		GivenName: claims.GivenName,
	}, nil
}
