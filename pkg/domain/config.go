package domain

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DevMode        bool   `default:"true" split_words:"true"`
	LogPretty      bool   `default:"true" split_words:"true"`
	LogPrettyColor bool   `default:"true" split_words:"true"`
	Addr           string `required:"true" split_words:"true"`

	PGUrl string `required:"true" split_words:"true"`

	NATSHostname                        string `split_words:"true"`
	NATSPort                            int    `split_words:"true"`
	NATSEmbeddedServer                  bool   `default:"true" split_words:"true"`
	NATSEmbeddedServerListen            bool   `split_words:"true"`
	NATSEmbeddedServerWebsocketHostname string `split_words:"true"`
	NATSEmbeddedServerWebsocketPort     int    `split_words:"true"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, Wrap(err, "cannot process configuration", nil)
}
