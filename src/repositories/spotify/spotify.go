package spotify

import (
	"context"

	"github.com/pkg/errors"
	sdk "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
	oauth2_spotify "golang.org/x/oauth2/spotify"

	"rollbringer/src"
	"rollbringer/src/domain/accounts"
)

type spotify struct {
	config      *src.Config
	oauthConfig *oauth2.Config
}

func New(config *src.Config) accounts.Spotify {
	return &spotify{
		config: config,
		oauthConfig: &oauth2.Config{
			Endpoint:     oauth2_spotify.Endpoint,
			ClientID:     config.SpotifyOauthClientID,
			ClientSecret: config.SpotifyOauthClientSecret,
			RedirectURL:  config.SpotifyOauthRedirectURL,
			Scopes:       []string{"user-read-private", "user-read-email"},
		},
	}
}

func (s *spotify) ConsentURL() (string, string) {
	state := src.CreateRandomString()
	return s.oauthConfig.AuthCodeURL(state), state
}

func (s *spotify) GetSpotifyUser(ctx context.Context, state string) (*accounts.SpotifyUser, error) {
	token, err := s.oauthConfig.Exchange(ctx, state)
	if err != nil {
		return nil, errors.Wrap(err, "cannot exchange state for token")
	}

	user, err := sdk.New(s.oauthConfig.Client(ctx, token)).CurrentUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get current user")
	}

	var profilePicture *string
	if user.Images != nil && len(user.Images) > 0 {
		profilePicture = &user.Images[0].URL
	}

	return &accounts.SpotifyUser{
		SpotifyID:      user.ID,
		DisplayName:    user.DisplayName,
		Email:          user.Email,
		ProfilePicture: profilePicture,
	}, nil
}
