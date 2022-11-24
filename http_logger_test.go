package forge_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/kyberbits/forge"
)

func TestLogger(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})

	logger := &forge.HTTPLogger{
		Logger: &forge.LoggerJSON{
			Encoder: json.NewEncoder(buffer),
		},
		Handler: &forge.HTTPRouter{
			NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("Logger test."))
			}),
		},
	}

	{ // Check root
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			panic(err)
		}
		testHandler(t, HandlerTestCase{
			Handler:            logger,
			Request:            request,
			ExpectedStatusCode: http.StatusCreated,
			ExpectedBody:       `Logger test.`,
		})
	}

	expected := forge.Log{
		Severity: "INFO",
		Message:  "HTTP Request",
		Context: map[string]interface{}{
			"handleTime":    0,
			"requestID":     "will override",
			"requestMethod": "GET",
			"requestPath":   "/",
			"requestQuery":  "",
			"responseCode":  http.StatusCreated,
		},
	}

	actual := forge.Log{}
	bufferBytes, _ := io.ReadAll(buffer)
	json.Unmarshal(bufferBytes, &actual)

	// Fix log
	expected.Context["requestID"] = "default"
	actual.Context["requestID"] = "default"
	expected.Timestamp = actual.Timestamp

	if err := forge.Assert(fmt.Sprint(expected), fmt.Sprint(actual)); err != nil {
		t.Error(err)
	}
}
