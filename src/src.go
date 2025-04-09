package src

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

var (
	ErrEntityConflict     = errors.New("entity conflict")
	ErrEntityNotFound     = errors.New("entity not found")
	ErrNoEntitiesEffected = errors.New("no entities effected")

	LevelTrace slog.Level = -8
	LevelDebug slog.Level = -4
	LevelInfo  slog.Level = 0
	LevelWarn  slog.Level = 4
	LevelError slog.Level = 8
	LevelFatal slog.Level = 12
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
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, errors.Wrap(err, "cannot process configuration")
}

func NewPrettyLogger(noColor bool) *slog.Logger {

	// Make logs pretty.
	h := tint.NewHandler(os.Stderr, &tint.Options{
		Level:   LevelTrace,
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
