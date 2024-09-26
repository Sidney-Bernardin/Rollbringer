//go:build all || pages
// +build all pages

package main

import (
	"net/http"

	"github.com/pkg/errors"

	handler "rollbringer/internal/handlers/pages"
	"rollbringer/internal/repositories/pubsub"
	"rollbringer/internal/services"
)

func init() {
	features["pages"] = func(deps globalDependencies) (http.Handler, services.BaseServicer, error) {
		return handler.NewHandler(deps.cfg, deps.logger), nil, nil
	}
}
