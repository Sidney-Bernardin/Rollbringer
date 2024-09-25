package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Port            string        `required:"true" split_words:"true"`
	ShutdownTimeout time.Duration `required:"false" default:"10s" split_words:"true"`
	CookiePath      string        `required:"true" split_words:"true"`

	UsersGoogleClientID     string        `required:"true" split_words:"true"`
	UsersGoogleClientSecret string        `required:"true" split_words:"true"`
	UsersRedirectURL        string        `required:"true" split_words:"true"`
	UsersSessionTimeout     time.Duration `required:"true" split_words:"true"`

	NATSHostname                     string        `required:"false" split_words:"true"`
	NATSPort                         int           `required:"false" split_words:"true"`
	NATSEmbeddedServer               bool          `required:"false" split_words:"true"`
	NATSEmbeddedServerListen         bool          `required:"false" split_words:"true"`
	NATSEmbeddedServerStartupTimeout time.Duration `required:"false" default:"10s" split_words:"true"`

	PostgresAddress string `required:"true" split_words:"true"`
}

func New() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, errors.Wrap(err, "cannot process configuration")
}
