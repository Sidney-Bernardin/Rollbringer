package service

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
)

func (svc *accountsService) NewGoogleUser(ctx context.Context, token *oauth2.Token) (*domain.User, error) {

	idTokenStr, ok := ctx.Value("token").(*oauth2.Token).Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token is not string")
	}

	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &GoogleIDTokenClaims{})
	if err != nil {
		return nil, domain.Wrap(err, "cannot parse ID token", nil)
	}
	claims := idToken.Claims.(*GoogleIDTokenClaims)

	return &domain.User{
		GoogleID:       &claims.Subject,
		Username:       claims.GivenName,
		ProfilePicture: claims.Picture,
		GoogleUser: &domain.GoogleUser{
			GoogleID:       claims.Subject,
			GivenName:      claims.GivenName,
			Email:          claims.Email,
			ProfilePicture: claims.Picture,
		},
	}, nil
}

func (svc *accountsService) NewSpotifyUser(ctx context.Context, oauthConfig *oauth2.Config, token *oauth2.Token) (*domain.User, error) {

	spotify := svc.spotifyRepo.Me(ctx, oauthConfig, token)
	spotifyPrivateUser, err := spotify.GetCurrentUser(ctx)
	if err != nil {
		return nil, domain.Wrap(err, "cannot get spotify current user", nil)
	}

	profilePicture := svc.Config.DefaultProfilePicture
	if spotifyPrivateUser.ProfilePicture != nil {
		profilePicture = *spotifyPrivateUser.ProfilePicture
	}

	return &domain.User{
		SpotifyID:      &spotifyPrivateUser.ID,
		Username:       spotifyPrivateUser.DisplayName,
		ProfilePicture: profilePicture,
		SpotifyUser: &domain.SpotifyUser{
			SpotifyID:      spotifyPrivateUser.ID,
			DisplayName:    spotifyPrivateUser.DisplayName,
			Email:          spotifyPrivateUser.Email,
			ProfilePicture: spotifyPrivateUser.ProfilePicture,
		},
	}, nil
}
