package service

import "github.com/golang-jwt/jwt/v5"

type GoogleIDTokenClaims struct {
	*jwt.RegisteredClaims

	GivenName string `json:"given_name"`
	Picture   string `json:"picture"`
}

type SpotifyUser struct {
	ID             string
	DisplayName    string
	ProfilePicture string
}

type SpotifyPrivateUser struct {
	SpotifyUser

	Email string
}
