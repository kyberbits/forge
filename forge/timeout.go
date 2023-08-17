package forge

import (
	"context"
	"errors"
	"net/http"
	"runtime"
	"time"

	"github.com/kyberbits/forge/forge/internal/stackparse"
)

// Replace http.TimeoutHandler with PanicReportTimeoutHandler.
func PanicReportTimeoutHandler(h http.Handler, dt time.Duration, logger *Logger) http.Handler {
	return http.TimeoutHandler(&panicReporterHandler{handler: h, logger: logger}, dt, "timeout")
}

type panicReporterHandler struct {
	handler http.Handler
	logger  *Logger
}

func (h *panicReporterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	defer func() {
		if panicValue := recover(); panicValue != nil {
			// Fix the stack trace
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			err, isErr := panicValue.(error)
			if !isErr {
				h.logger.Error(r.Context(), "HTTP Panic (Value)", map[string]any{
					"value":    panicValue,
					"method":   r.Method,
					"path":     r.URL.Path,
					"duration": uint64(time.Since(startTime).Seconds()),
					"stack":    stackparse.Parse(buf),
				})
			}

			if errors.Is(err, context.Canceled) {
				h.logger.Warning(r.Context(), "Context Canceled", map[string]any{
					"err":      err,
					"method":   r.Method,
					"path":     r.URL.Path,
					"duration": uint64(time.Since(startTime).Seconds()),
				})

				return
			}

			if errors.Is(err, http.ErrAbortHandler) {
				h.logger.Error(r.Context(), "HTTP Timeout", map[string]any{
					"err":      err,
					"method":   r.Method,
					"path":     r.URL.Path,
					"duration": uint64(time.Since(startTime).Seconds()),
					"stack":    stackparse.Parse(buf),
				})

				return
			}

			h.logger.Error(r.Context(), "HTTP Panic (Error)", map[string]any{
				"err":      err,
				"method":   r.Method,
				"path":     r.URL.Path,
				"duration": uint64(time.Since(startTime).Seconds()),
				"stack":    stackparse.Parse(buf),
			})
		}
	}()
	h.handler.ServeHTTP(w, r)
}
