package forge

import (
	"net/http"
	"time"
)

type HTTPLogger struct {
	Logger  *Logger
	Handler http.Handler
}

func (httpLogger *HTTPLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recorder := &statusRecorder{
		ResponseWriter: w,
		ResponseCode:   200,
	}

	startTime := time.Now()

	httpLogger.Handler.ServeHTTP(recorder, r)

	duration := time.Since(startTime)

	// Log the completed request
	httpLogger.Logger.Info(
		r.Context(),
		"HTTP Request",
		map[string]interface{}{
			"requestAddr":   getRemoteAddr(r),
			"requestMethod": r.Method,
			"requestPath":   r.URL.Path,
			"requestQuery":  r.URL.Query().Encode(),
			"responseCode":  recorder.ResponseCode,
			"handleTime":    duration.Milliseconds(),
		},
	)
}

type statusRecorder struct {
	http.ResponseWriter
	ResponseCode int
}

func (recorder *statusRecorder) WriteHeader(status int) {
	recorder.ResponseCode = status
	recorder.ResponseWriter.WriteHeader(status)
}

func getRemoteAddr(r *http.Request) string {
	// The X-FORWARDED-FOR header is the de-facto standard header for identifying the originating IP address of a client connecting
	// to a web server through an HTTP proxy or a load balancer. https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	if r.Header.Get("X-Forwarded-For") != "" {
		return r.Header.Get("X-Forwarded-For")
	}

	return r.RemoteAddr
}
