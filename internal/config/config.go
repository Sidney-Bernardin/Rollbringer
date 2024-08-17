package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `required:"true" split_words:"true"`

	GoogleClientID     string `required:"true" split_words:"true"`
	GoogleClientSecret string `required:"true" split_words:"true"`

	RedirectURL        string        `required:"true" split_words:"true"`
	UserSessionTimeout time.Duration `required:"true" split_words:"true"`

	PostgresAddress string `required:"true" split_words:"true"`

	NATSEmbeddedServer               bool          `required:"true" split_words:"true"`
	NATSListenWithEmbeddedServer     bool          `required:"true" split_words:"true"`
	NATSEmbeddedServerStartupTimeout time.Duration `required:"true" split_words:"true"`
	NATSHostname                     string        `required:"true" split_words:"true"`
	NATSPort                         int           `required:"true" split_words:"true"`
}

func New() (cfg *Config, err error) {
	err = envconfig.Process("APP", &cfg)
	return cfg, err
}
