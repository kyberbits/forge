package forge_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/kyberbits/forge/forge"
	"github.com/kyberbits/forge/forgetest"
)

func TestHTTPClientSuccess(t *testing.T) {
	client := forge.NewHTTPClient(
		&http.Client{
			Transport: forgetest.MockRoundTripperQueue(t, []forgetest.MockRoundTripFunc{
				func(t *testing.T, request *http.Request) (*http.Response, error) {
					username, password, _ := request.BasicAuth()
					if err := forgetest.Assert("username", username); err != nil {
						t.Fatal(err)
					}

					if err := forgetest.Assert("password", password); err != nil {
						t.Fatal(err)
					}

					file, _ := os.Open("test_files/client/test.json")

					return &http.Response{
						Request:    request,
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(file),
					}, nil
				},
			}),
		},
		[]forge.HTTPClientMiddleware{
			forge.HTTPBasicAuth{
				Username: "username",
				Password: "password",
			},
		},
	)

	type Target struct {
		Greeting string `json:"greeting"`
	}

	target := Target{}
	if err := client.JSONRequest(context.Background(), http.MethodGet, "/", nil, &target); err != nil {
		t.Fatal(err)
	}

	expected := Target{
		Greeting: "Hello there.",
	}
	actual := target

	if err := forgetest.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}
