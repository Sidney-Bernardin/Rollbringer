package service

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
)

func (svc *accountsService) LoginWithGoogle(ctx context.Context, token *oauth2.Token) (*domain.User, error) {

	idTokenStr, ok := ctx.Value("token").(*oauth2.Token).Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token is not string")
	}

	idToken, _, err := jwt.NewParser().ParseUnverified(idTokenStr, &GoogleIDTokenClaims{})
	if err != nil {
		return nil, domain.Wrap(err, "cannot parse ID token", nil)
	}
	claims := idToken.Claims.(*GoogleIDTokenClaims)

	user := &domain.User{
		Username: claims.GivenName,
		GoogleUser: &domain.GoogleUser{
			GoogleID:       claims.Subject,
			GivenName:      claims.GivenName,
			ProfilePicture: claims.Picture,
		},
	}

	if err := svc.login(ctx, user); err != nil {
		return nil, domain.Wrap(err, "cannot login", nil)
	}

	return user, nil
}

func (svc *accountsService) LoginWithSpotify(ctx context.Context, oauthConfig *oauth2.Config, token *oauth2.Token) (*domain.User, error) {

	spotify := svc.spotifyRepo.Me(ctx, oauthConfig, token)
	spotifyPrivateUser, err := spotify.GetCurrentUser(ctx)
	if err != nil {
		return nil, domain.Wrap(err, "cannot get spotify current user", nil)
	}

	user := &domain.User{
		Username: spotifyPrivateUser.DisplayName,
		SpotifyUser: &domain.SpotifyUser{
			SpotifyID:      spotifyPrivateUser.ID,
			DisplayName:    spotifyPrivateUser.DisplayName,
			ProfilePicture: spotifyPrivateUser.ProfilePicture,
		},
	}

	if err := svc.login(ctx, user); err != nil {
		return nil, domain.Wrap(err, "cannot login", nil)
	}

	return user, nil
}

func (svc *accountsService) login(ctx context.Context, user *domain.User) error {

	err := svc.accountsDBRepo.Transaction(ctx, func(tx AccountsDatabaseRepository) error {

		if err := tx.UserInsert(ctx, user); err != nil {
			return domain.Wrap(err, "cannot insert user", nil)
		}

		csrfToken, err := domain.NewRandomString(ctx)
		if err != nil {
			return domain.Wrap(err, "cannot create CSRF token", nil)
		}

		user.Session = &domain.Session{
			UserID:    user.ID,
			CSRFToken: csrfToken,
		}

		if err := tx.SessionInsert(ctx, user.Session); err != nil {
			return domain.Wrap(err, "cannot insert session", nil)
		}

		if user.GoogleUser != nil {
			user.GoogleUser.UserID = user.ID

			if err := tx.GoogleUserInsert(ctx, user.GoogleUser); err != nil {
				return domain.Wrap(err, "cannot insert google-user", nil)
			}
		}

		if user.SpotifyUser != nil {
			user.SpotifyUser.UserID = user.ID

			if err := tx.SpotifyUserInsert(ctx, user.SpotifyUser); err != nil {
				return domain.Wrap(err, "cannot insert spotify-user", nil)
			}
		}

		return nil
	})

	return domain.Wrap(err, "cannot do transaction", nil)
}
