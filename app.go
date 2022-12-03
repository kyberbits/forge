package forge

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type App interface {
	// The Logger to be used
	Logger() *Logger
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

func KeepRunning(ctx context.Context, app App, action func(ctx context.Context), coolDown time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				app.Logger().Error(ctx, "Panic Err", map[string]any{
					"err": err,
				})
			} else {
				app.Logger().Error(ctx, "Panic Err", map[string]any{
					"err": r,
				})
			}

			select {
			case <-time.After(coolDown):
				KeepRunning(ctx, app, action, coolDown)
			case <-ctx.Done():
				return
			}
		}
	}()

	action(ctx)

	select {
	case <-time.After(coolDown):
		KeepRunning(ctx, app, action, coolDown)
	case <-ctx.Done():
		return
	}

}

func Run(ctx context.Context, app WebApp) error {
	go KeepRunning(ctx, app, app.Background, time.Second*5)

	httpServer := &http.Server{
		Addr:         app.ListenAddress(),
		Handler:      http.TimeoutHandler(app.Handler(), time.Second*30, "Timeout"),
		ErrorLog:     app.Logger().StandardLogger(),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	app.Logger().Info(ctx, "Serving web application", map[string]interface{}{
		"listen": fmt.Sprintf("http://%s", app.ListenAddress()),
	})

	if err := httpServer.ListenAndServe(); err != nil {
		app.Logger().Critical(ctx, "Webserver failed to start", map[string]interface{}{
			"addr": app.ListenAddress(),
			"err":  err,
		})
		return err
	}

	return nil
}
