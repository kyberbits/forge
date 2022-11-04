package forge_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/kyberbits/forge"
)

type TestInterceptor struct {
	ModRequestFunc func(request *http.Request) error
	ErrorCheckFunc func(response *http.Response) error
}

func (i *TestInterceptor) ModRequest(request *http.Request) error {
	if i.ModRequestFunc == nil {
		return nil
	}

	return i.ModRequestFunc(request)
}

func (i *TestInterceptor) ErrorCheck(response *http.Response) error {
	if i.ErrorCheckFunc == nil {
		return nil
	}

	return i.ErrorCheckFunc(response)
}

func TestHTTPClientSuccess(t *testing.T) {
	client := forge.NewHTTPClient(
		&TestInterceptor{},
		&http.Client{
			Transport: forge.MockRoundTripperQueue(t, []forge.MockRoundTripFunc{
				func(t *testing.T, request *http.Request) (*http.Response, error) {
					file, _ := os.Open("test_files/client/test.json")
					return &http.Response{
						Request:    request,
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(file),
					}, nil
				},
			}),
		},
	)

	type Target struct {
		Greeting string `json:"greeting"`
	}

	target := Target{}
	if err := client.Request(http.MethodGet, "/", nil, &target); err != nil {
		t.Fatal(err)
	}

	expected := Target{
		Greeting: "Hello there.",
	}
	actual := target
	if err := forge.Assert(expected, actual); err != nil {
		t.Fatal(err)
	}
}
