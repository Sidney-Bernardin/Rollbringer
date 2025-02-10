package domain

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr                  string `required:"true" split_words:"true"`
	DevMode               bool   `default:"true" split_words:"true"`
	DefaultProfilePicture string `default:"/pages/static/favicon.png" split_words:"true"`

	LogPretty      bool `default:"true" split_words:"true"`
	LogPrettyColor bool `default:"true" split_words:"true"`

	OauthGoogleClientID          string `required:"true" split_words:"true"`
	OauthGoogleClientSecret      string `required:"true" split_words:"true"`
	OauthGoogleSigninRedirectURL string `required:"true" split_words:"true"`
	OauthGoogleSignupRedirectURL string `required:"true" split_words:"true"`

	OauthSpotifyClientID          string `required:"true" split_words:"true"`
	OauthSpotifyClientSecret      string `required:"true" split_words:"true"`
	OauthSpotifySigninRedirectURL string `required:"true" split_words:"true"`
	OauthSpotifySignupRedirectURL string `required:"true" split_words:"true"`

	OAuthStateCookieTimeout time.Duration `default:"2m" split_words:"true"`

	SessionCookieTimeout time.Duration `required:"true" split_words:"true"`
	SessionCookiePath    string        `required:"true" split_words:"true"`

	PGUrl string `required:"true" split_words:"true"`

	NATSHostname             string `split_words:"true"`
	NATSPort                 int    `split_words:"true"`
	NATSEmbeddedServer       bool   `default:"true" split_words:"true"`
	NATSEmbeddedServerListen bool   `split_words:"true"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, Wrap(err, "cannot process configuration", nil)
}
