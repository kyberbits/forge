package forge_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyberbits/forge/forge"
	"github.com/kyberbits/forge/forgetest"
)

func TestLoggerMiddleware(t *testing.T) {
	ctx := context.Background()
	buffer := bytes.NewBuffer([]byte{})

	httpLogger := &forge.HTTPLogger{
		Logger: forge.NewLogger(
			"http",
			buffer,
			func(logEntry *forge.Log, ctx context.Context, r *http.Request) {
			},
		),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("Logger test."))
		}),
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	httpLogger.ServeHTTP(recorder, request)

	actual := forge.Log{}
	bufferBytes, _ := io.ReadAll(buffer)

	if err := json.Unmarshal(bufferBytes, &actual); err != nil {
		t.Fatal(err)
	}

	if err := forgetest.Assert("HTTP Request", actual.Message); err != nil {
		t.Fatal(actual)
	}
}
