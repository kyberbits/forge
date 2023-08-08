package forge

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/aaronellington/environment-go/environment"
)

func NewRuntime() *Runtime {
	return &Runtime{
		Environment: environment.New(),
		FS:          os.DirFS("."),
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Stdin:       os.Stdin,
	}
}

type Runtime struct {
	Environment *environment.Environment
	Stdout      io.Writer
	Stderr      io.Writer
	Stdin       io.Reader
	FS          fs.FS
}

func (runtime *Runtime) Execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = runtime.Stderr
	cmd.Stdin = runtime.Stdin
	cmd.Stdout = runtime.Stdout

	return cmd.Run()
}

func (runtime *Runtime) KeepRunning(ctx context.Context, app App, action func(ctx context.Context), coolDown time.Duration) {
	ticker := time.NewTicker(coolDown)

	tickFunc := func() {
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
			}
		}()
		action(ctx)
	}

	tickFunc()

	for range ticker.C {
		tickFunc()
	}
}

func (runtime *Runtime) Serve(ctx context.Context, app App) error {
	go runtime.KeepRunning(ctx, app, app.Background, time.Second*5)

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
