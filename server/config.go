package server

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	HttpAddr                    string        `required:"true" split_words:"true"`
	HttpGracfullShutdownTimeout time.Duration `default:"3s" split_words:"true"`

	SessionTimeout time.Duration `default:"12h" split_words:"true"`

	PostgresUrl    string   `required:"true" split_words:"true"`
	CassandraHosts []string `required:"true" split_words:"true"`
	RedisAddr      string   `required:"true" split_words:"true"`
	RedisPassword  string   `required:"true" split_words:"true"`
	NatsUrl        string   `required:"true" split_words:"true"`

	GoogleOauthClientId     string `required:"true" split_words:"true"`
	GoogleOauthClientSecret string `required:"true" split_words:"true"`
	GoogleOauthRedirectUrl  string `required:"true" split_words:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, errors.WithStack(err)
}
