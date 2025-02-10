package spotify

import (
	"context"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"

	"rollbringer/pkg/domain"
	service "rollbringer/pkg/domain/services/accounts"
)

type spotifyRepository struct {
	client *spotify.Client
}

func NewSpotifyRepository() service.SpotifyRepository {
	return &spotifyRepository{}
}

func (repo *spotifyRepository) Me(ctx context.Context, oauthConfig *oauth2.Config, token *oauth2.Token) service.SpotifyRepository {
	return &spotifyRepository{
		client: spotify.New(oauthConfig.Client(ctx, token)),
	}
}

func (repo *spotifyRepository) GetCurrentUser(ctx context.Context) (ret *service.SpotifyPrivateUser, err error) {
	privateUser, err := repo.client.CurrentUser(ctx)
	if err != nil {
		return nil, domain.Wrap(err, "cannot get current user", nil)
	}

	ret = &service.SpotifyPrivateUser{
		SpotifyUser: service.SpotifyUser{
			ID:          privateUser.ID,
			DisplayName: privateUser.DisplayName,
		},
		Email: privateUser.Email,
	}

	if privateUser.Images != nil && len(privateUser.Images) > 0 {
		ret.ProfilePicture = &privateUser.Images[0].URL
	}

	return ret, nil
}
