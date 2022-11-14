package forge

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type App interface {
	// The Logger to be used
	Logger() Logger
}

type WebApp interface {
	App
	// How it responds to web requests
	Handler() http.Handler
	// What address the webserver should listen on
	ListenAddress() string
	// Tasks to run in the background (a Go routine)
	Background(ctx context.Context)
}

func Run(ctx context.Context, app WebApp) error {
	go app.Background(ctx)

	httpServer := &http.Server{
		Addr:         app.ListenAddress(),
		Handler:      http.TimeoutHandler(app.Handler(), time.Second*30, "Timeout"),
		ErrorLog:     app.Logger().StandardLogger(),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	app.Logger().Info("Serving web application", map[string]interface{}{
		"listen": fmt.Sprintf("http://%s", app.ListenAddress()),
	})

	if err := httpServer.ListenAndServe(); err != nil {
		app.Logger().Critical("Webserver failed to start", map[string]interface{}{
			"addr": app.ListenAddress(),
			"err":  err,
		})
		return err
	}

	return nil
}
