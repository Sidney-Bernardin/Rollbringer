package service

import (
	"context"

	"rollbringer/pkg/domain"

	"github.com/pkg/errors"
)

func (svc *accountsService) Signup(ctx context.Context, user *domain.User) error {
	err := svc.accountsDBRepo.Transaction(ctx, func(tx AccountsDatabaseRepository) (err error) {

		if user.GoogleUser != nil {
			if err := tx.GoogleUserInsert(ctx, user.GoogleUser); err != nil {
				if errors.Cause(err) == domain.ErrAlreadyExists {
					return domain.UserErr(ctx, domain.UsrErrTypeGoogleUserAlreadyExists, "This Google account is already linked with a Rollbringer account.", nil)
				}

				return domain.Wrap(err, "cannot insert google-user", nil)
			}
		}

		if user.SpotifyUser != nil {
			if err := tx.SpotifyUserInsert(ctx, user.SpotifyUser); err != nil {
				if errors.Cause(err) == domain.ErrAlreadyExists {
					return domain.UserErr(ctx, domain.UsrErrTypeSpotifyUserAlreadyExists, "This Spotify account is already linked with a Rollbringer account.", nil)
				}

				return domain.Wrap(err, "cannot insert spotify-user", nil)
			}
		}

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

		return nil
	})

	return domain.Wrap(err, "cannot do transaction", nil)
}

func (svc *accountsService) Signin(ctx context.Context, user *domain.User) error {
	err := svc.accountsDBRepo.Transaction(ctx, func(tx AccountsDatabaseRepository) error {

		if user.GoogleUser != nil {
			err := tx.GoogleUserUpdate(ctx, "google_id", user.GoogleUser.GoogleID, map[string]any{
				"given_name":      user.GoogleUser.GivenName,
				"email":           user.GoogleUser.Email,
				"profile_picture": user.GoogleUser.ProfilePicture,
			})

			if err != nil {
				if errors.Cause(err) == domain.ErrNotFound {
					return domain.UserErr(ctx, domain.UsrErrTypeGoogleUserDoesNotExists, "This Google account is not linked with a Rollbringer account.", nil)
				}

				return domain.Wrap(err, "cannot update google-user", nil)
			}

			u, err := tx.UserGet(ctx, "google_id", user.GoogleID)
			if err != nil {
				return domain.Wrap(err, "cannot get user by google-id", nil)
			}
			*user = *u
		}

		if user.SpotifyUser != nil {
			err := tx.SpotifyUserUpdate(ctx, "spotify_id", user.SpotifyUser.SpotifyID, map[string]any{
				"display_name":    user.SpotifyUser.DisplayName,
				"email":           user.SpotifyUser.Email,
				"profile_picture": user.SpotifyUser.ProfilePicture,
			})

			if err != nil {
				if errors.Cause(err) == domain.ErrNotFound {
					return domain.UserErr(ctx, domain.UsrErrTypeSpotifyUserDoesNotExists, "This Spotify account is not linked with a Rollbringer account.", nil)
				}

				return domain.Wrap(err, "cannot update spotify-user", nil)
			}

			u, err := tx.UserGet(ctx, "spotify_id", user.SpotifyID)
			if err != nil {
				return domain.Wrap(err, "cannot get user by spotify-id", nil)
			}
			*user = *u
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

		return nil
	})

	return domain.Wrap(err, "cannot do transaction", nil)
}
