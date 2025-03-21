package src

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	DevMode bool `default:"true" split_words:"true"`

	GracfullShutdownTimeout time.Duration `default:"3s" split_words:"true"`
	APIAddr                 string        `required:"true" split_words:"true"`

	PostgresURL string `required:"true" split_words:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, errors.Wrap(err, "cannot process configuration")
}
