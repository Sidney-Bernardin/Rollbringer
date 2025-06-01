package google

import (
	"context"

	"github.com/Sidney-Bernardin/Rollbringer/server"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	*jwt.RegisteredClaims

	GivenName string `json:"given_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
}

type Google struct {
	config *server.Config

	oauthConfig *oauth2.Config
}

func New(config *server.Config) *Google {
	return &Google{
		config: config,
		oauthConfig: &oauth2.Config{
			Endpoint:     google.Endpoint,
			ClientID:     config.GoogleOauthClientId,
			ClientSecret: config.GoogleOauthClientSecret,
			RedirectURL:  config.GoogleOauthRedirectUrl,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}
}

func (g *Google) ConsentURL() (url string, state string) {
	state = server.CreateRandomString()
	return g.oauthConfig.AuthCodeURL(state), state
}

func (g *Google) GetGoogleUser(ctx context.Context, code string) (*GoogleUser, error) {

	token, err := g.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, &server.UserError{Type: server.UserErrorTypeUnauthorized}
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token is not a string")
	}

	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &GoogleUser{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse ID-token")
	}

	return idToken.Claims.(*GoogleUser), nil
}
