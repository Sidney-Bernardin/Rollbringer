package google

import (
	"context"

	"github.com/Sidney-Bernardin/Rollbringer/server"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type GoogleUser struct {
	*jwt.RegisteredClaims

	GivenName string `json:"given_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
}

type Google struct {
	OAuthConfig *oauth2.Config
}

func (g *Google) GetGoogleUser(ctx context.Context, code string) (*GoogleUser, error) {

	token, err := g.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, server.NewUserError(server.UserErrorTypeUnauthorized, "", nil)
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
