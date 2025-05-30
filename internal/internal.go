package internal

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"
	"time"

	"github.com/dpotapov/slogpfx"
	"github.com/kelseyhightower/envconfig"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

type Config struct {
	DevMode                 bool          `default:"true" split_words:"true"`
	GracfullShutdownTimeout time.Duration `default:"3s" split_words:"true"`

	ServerAddr string `required:"true" split_words:"true"`

	NatsURL        string   `required:"true" split_words:"true"`
	PostgresURL    string   `required:"true" split_words:"true"`
	CassandraHosts []string `required:"true" split_words:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, errors.Wrap(err, "cannot process configuration")
}

func NewPrettyLogger(noColor bool) *slog.Logger {

	// Make logs pretty.
	h := tint.NewHandler(os.Stderr, &tint.Options{
		NoColor: noColor && !isatty.IsTerminal(os.Stderr.Fd()),
	})

	// Prominently displays "namespace" attributes.
	h = slogpfx.NewHandler(h, &slogpfx.HandlerOptions{
		PrefixKeys: []string{"namespace"},
	})

	return slog.New(h)
}

func CreateRandomString() string {
	var bState = make([]byte, 32)
	if _, err := rand.Read(bState); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bState)
}
