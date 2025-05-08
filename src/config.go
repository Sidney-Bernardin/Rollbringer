package src

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	DevMode                 bool          `default:"true" split_words:"true"`
	GracfullShutdownTimeout time.Duration `default:"3s" split_words:"true"`

	APIAddr string `required:"true" split_words:"true"`

	SessionCookieTimeout time.Duration `required:"true" split_words:"true"`
	OAuthCookieTimeout   time.Duration `default:"15m" split_words:"true"`

	GoogleOauthClientID     string `required:"true" split_words:"true"`
	GoogleOauthClientSecret string `required:"true" split_words:"true"`
	GoogleOauthRedirectURL  string `required:"true" split_words:"true"`

	SpotifyOauthClientID     string `required:"true" split_words:"true"`
	SpotifyOauthClientSecret string `required:"true" split_words:"true"`
	SpotifyOauthRedirectURL  string `required:"true" split_words:"true"`

	PostgresAccountsURL string `required:"true" split_words:"true"`
	PostgresPlayURL     string `required:"true" split_words:"true"`

	NatsURL string `required:"true" split_words:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, errors.Wrap(err, "cannot process configuration")
}
