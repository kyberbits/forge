package forge

import (
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
	defer func() {
		if err := recover(); err != nil {
			if err == http.ErrAbortHandler {
				h.logger.Warning(r.Context(), "HTTP Timeout", map[string]any{
					"path": r.URL.Path,
					"err":  err,
				})

				return
			}

			// debug.Stack()

			// // Fix the stack trace
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			h.logger.Error(r.Context(), "HTTP Panic", map[string]any{
				"err":   err,
				"stack": stackparse.Parse(buf),
			})
		}
	}()
	h.handler.ServeHTTP(w, r)
}
