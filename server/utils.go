package server

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/dpotapov/slogpfx"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

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
