package forge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func NewRuntime() *Runtime {
	return &Runtime{
		Environment: NewEnvironment(),
		FS:          os.DirFS("."),
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Stdin:       os.Stdin,
	}
}

type Runtime struct {
	Environment Environment
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

func (runtime *Runtime) ReadInDefaultEnvironmentFiles() error {
	defaultFiles := []string{
		// Values already set in the Environment will not be changed
		".env.local", // Not tracked in git, first priority
		".env",       // Defaults if not set other
	}

	for _, defaultFile := range defaultFiles {
		if err := runtime.ReadInEnvironmentFile(defaultFile); err != nil {
			return err
		}
	}

	return nil
}

func (runtime *Runtime) ReadInEnvironmentFile(fileName string) error {
	file, err := runtime.FS.Open(fileName)
	if err != nil {
		// File does not exist errors are okay
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	fileContentBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	fileContents := string(fileContentBytes)

	return runtime.Environment.ImportEnvFileContents(fileContents)
}

func (runtime *Runtime) KeepRunning(ctx context.Context, app App, action func(ctx context.Context), coolDown time.Duration) {
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
				runtime.KeepRunning(ctx, app, action, coolDown)
			case <-ctx.Done():
				return
			}
		}
	}()

	action(ctx)

	select {
	case <-time.After(coolDown):
		runtime.KeepRunning(ctx, app, action, coolDown)
	case <-ctx.Done():
		return
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
