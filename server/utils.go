package server

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/dpotapov/slogpfx"
	"github.com/lmittmann/tint"
)

func NewPrettyLogger() *slog.Logger {

	// Make logs pretty.
	h := tint.NewHandler(os.Stderr, nil)

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
