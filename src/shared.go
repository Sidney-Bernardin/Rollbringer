package src

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

	APIAddr string `required:"true" split_words:"true"`

	SessionCookieTimeout    time.Duration `required:"true" split_words:"true"`
	SessionCookiePath       string        `required:"true" split_words:"true"`
	OAuthStateCookieTimeout time.Duration `default:"2m" split_words:"true"`

	OauthGoogleClientID     string `required:"true" split_words:"true"`
	OauthGoogleClientSecret string `required:"true" split_words:"true"`
	OauthGoogleRedirectURL  string `required:"true" split_words:"true"`

	OauthSpotifyClientID     string `required:"true" split_words:"true"`
	OauthSpotifyClientSecret string `required:"true" split_words:"true"`
	OauthSpotifyRedirectURL  string `required:"true" split_words:"true"`

	PostgresAccountsURL string `required:"true" split_words:"true"`
	PostgresPlayURL     string `required:"true" split_words:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	return &cfg, errors.Wrap(err, "cannot process configuration")
}

/////

var (
	LevelTrace slog.Level = -8
	LevelDebug slog.Level = -4
	LevelInfo  slog.Level = 0
	LevelWarn  slog.Level = 4
	LevelError slog.Level = 8
	LevelFatal slog.Level = 12
)

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

/////

type ExternalErrorType string

const (
	ExternalErrorTypeUUIDInvalid    = "uuid_invalid"
	ExternalErrorTypeEntityNotFound = "entity_not_found"
)

type ExternalError struct {
	Type        ExternalErrorType
	Description string
	Attrs       map[string]any
}

func (err *ExternalError) Error() string {
	return fmt.Sprintf("%s: %s", err.Type, err.Description)
}

/////

func NewRandomString(ctx context.Context) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "cannot read bytes")
	}
	return hex.EncodeToString(b), nil
}
