package models

import (
	"errors"
	"rollbringer/src"
)

type User struct {
	ID src.UUID

	GoogleID   *string
	GoogleUser *GoogleUser

	SpotifyID   *string
	SpotifyUser *SpotifyUser

	Username       Username
	ProfilePicture string
}

func NewUser(googleUser *GoogleUser, spotifyUser *SpotifyUser) (user *User, err error) {
	user = &User{
		ID:             src.NewUUID(),
		ProfilePicture: "/static/favicon.png",
	}

	if googleUser != nil {
		user.GoogleID = &googleUser.GoogleID
		user.GoogleUser = googleUser
		user.Username = Username(googleUser.GivenName)
		user.ProfilePicture = googleUser.ProfilePicture
	} else if spotifyUser != nil {
		user.SpotifyID = &spotifyUser.SpotifyID
		user.SpotifyUser = spotifyUser
		user.Username = Username(spotifyUser.DisplayName)
		if spotifyUser.ProfilePicture != nil {
			user.ProfilePicture = *spotifyUser.ProfilePicture
		}
	} else {
		return nil, errors.New("user must have a provider")
	}

	return user, nil
}

type GoogleUser struct {
	GoogleID string

	GivenName      string
	Email          string
	ProfilePicture string
}

type SpotifyUser struct {
	SpotifyID string

	DisplayName    string
	Email          string
	ProfilePicture *string
}

type Username string

const ExternalErrorTypeInvalidUsername src.ExternalErrorType = "invalid_username"

func ParseUsername(str string) (Username, error) {
	if len(str) == 0 || 25 < len(str) {
		return "", &src.ExternalError{
			Type:    ExternalErrorTypeInvalidUsername,
			Msg:     "Must be between 1 and 25 characters",
			Details: map[string]any{"username": str},
		}
	}

	return Username(str), nil
}
