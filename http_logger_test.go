package forge_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyberbits/forge"
)

func TestLoggerMiddleware(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})

	httpLogger := &forge.HTTPLogger{
		Logger: *forge.NewLogger(
			"http",
			buffer,
			func(logEntry *forge.LogEntry, ctx context.Context, r *http.Request) {
			},
		),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Logger test."))
		}),
	}

	request, err := http.NewRequest(http.MethodGet, "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	httpLogger.ServeHTTP(recorder, request)

	actual := forge.LogEntry{}
	bufferBytes, _ := io.ReadAll(buffer)
	json.Unmarshal(bufferBytes, &actual)

	if err := forge.Assert("HTTP Request", actual.Message); err != nil {
		t.Fatal(actual)
	}
}
