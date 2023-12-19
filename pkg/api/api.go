package api

import (
	database "rollbringer/pkg/repositories/database"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type API struct {
	DB     *database.Database
	Logger *zerolog.Logger

	GoogleOAuthConfig *oauth2.Config
}
