package forge

import (
	"context"
	"net/http"
)

type App interface {
	// The Logger to be used
	Logger() *Logger
	// How it responds to web requests
	Handler() http.Handler
	// What address the webserver should listen on
	ListenAddress() string
	// Tasks to run in the background (a Go routine)
	Background(ctx context.Context)
}
