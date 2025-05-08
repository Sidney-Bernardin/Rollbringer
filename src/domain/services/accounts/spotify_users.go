package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src/domain"
)

type SpotifyUser struct {
	SpotifyID      string
	DisplayName    string
	Email          string
	ProfilePicture *string
}

func (svc *service) SpotifyLogin(ctx context.Context, oauthCode string, newAccount bool) (uuid.UUID, error) {

	spotifyUser, err := svc.spotify.GetSpotifyUser(ctx, oauthCode)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "spotify cannot get spotify-user")
	}

	if !newAccount {

		sessionID, err := svc.database.SpotifySignin(ctx, spotifyUser)
		if errors.Is(err, domain.ErrNoEntitiesEffected) {
			return uuid.Nil, &domain.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Spotify account is not linked with a Rollbringer account."}
		}

		return sessionID, errors.Wrap(err, "database cannot signin")
	}

	user := &User{
		ID:             uuid.New(),
		SpotifyID:      &spotifyUser.SpotifyID,
		Username:       Username(spotifyUser.DisplayName),
		ProfilePicture: "/static/favicon.png",
	}

	if spotifyUser.ProfilePicture != nil {
		user.ProfilePicture = *spotifyUser.ProfilePicture
	}

	sessionID, err := svc.database.SpotifySignup(ctx, spotifyUser, user)
	if errors.Is(err, domain.ErrEntityConflict) {
		return uuid.Nil, &domain.ExternalError{Type: ExternalErrorTypeProviderNotLinked, Msg: "The Spotify account is already linked with a Rollbringer account."}
	}

	return sessionID, errors.Wrap(err, "database cannot signup")
}
