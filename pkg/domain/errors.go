package domain

import "github.com/pkg/errors"

var (
	ErrUnauthorized = errors.New("unauthorized")

	ErrUserNotFound = errors.New("user was not found")

	ErrMaxGames     = errors.New("max games reached")
	ErrGameNotFound = errors.New("game was not found")

	ErrPlayMaterialNotFound = errors.New("play material was not found")
	ErrInvalidPDFName       = errors.New("invalid pdf name")
)
