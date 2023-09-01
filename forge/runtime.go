package forge

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/aaronellington/environment-go/environment"
	"github.com/kyberbits/forge/forge/internal/stackparse"
	"github.com/kyberbits/forge/forgeutils"
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

func (r *Runtime) Execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = r.Stderr
	cmd.Stdin = r.Stdin
	cmd.Stdout = r.Stdout

	return cmd.Run()
}

func (r *Runtime) KeepRunning(ctx context.Context, app App, action func(ctx context.Context), coolDown time.Duration) {
	ticker := time.NewTicker(coolDown)

	tickFunc := func() {
		defer func() {
			if panicValue := recover(); panicValue != nil {
				// Fix the stack trace
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]

				if err, ok := panicValue.(error); ok {
					app.Logger().Error(ctx, "Panic - KeepRunning (Err)", map[string]any{
						"err":   err,
						"stack": stackparse.Parse(buf),
					})
				} else {
					app.Logger().Error(ctx, "Panic - KeepRunning (Value)", map[string]any{
						"value": panicValue,
						"stack": stackparse.Parse(buf),
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

func (r *Runtime) Serve(ctx context.Context, app App) error {
	go r.KeepRunning(ctx, app, app.Background, time.Second*5)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = forgeutils.ContextAddToRequest(r)
		PanicReportTimeoutHandler(app.Handler(), time.Second*30, app.Logger()).ServeHTTP(w, r)
	})

	httpServer := &http.Server{
		Addr:         app.ListenAddress(),
		Handler:      handler,
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
