package google

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	oauth2_google "golang.org/x/oauth2/google"

	"rollbringer/src"
	"rollbringer/src/domain/accounts"
)

type googleIDTokenClaims struct {
	*jwt.RegisteredClaims

	GivenName string `json:"given_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
}

type google struct {
	config      *src.Config
	oauthConfig *oauth2.Config
}

func New(config *src.Config) accounts.Google {
	return &google{
		config: config,
		oauthConfig: &oauth2.Config{
			Endpoint:     oauth2_google.Endpoint,
			ClientID:     config.GoogleOauthClientID,
			ClientSecret: config.GoogleOauthClientSecret,
			RedirectURL:  config.GoogleOauthRedirectURL,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}
}

func (g *google) ConsentURL() (string, string) {
	state := src.CreateRandomString()
	return state, g.oauthConfig.AuthCodeURL(state)
}

func (g *google) GetGoogleUser(ctx context.Context, state string) (*accounts.GoogleUser, error) {
	token, err := g.oauthConfig.Exchange(ctx, state)
	if err != nil {
		return nil, errors.Wrap(err, "cannot exchange state for token")
	}

	idTokenStr, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token is not string")
	}

	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &googleIDTokenClaims{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse ID-token")
	}
	claims := idToken.Claims.(*googleIDTokenClaims)

	return &accounts.GoogleUser{
		GoogleID:       claims.Subject,
		GivenName:      claims.GivenName,
		Email:          claims.Email,
		ProfilePicture: claims.Picture,
	}, nil
}
