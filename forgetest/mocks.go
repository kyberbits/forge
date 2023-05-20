package forgetest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

type MockableRoundTripper struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (mockable MockableRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return mockable.RoundTripFunc(request)
}

type MockRoundTripFunc func(t *testing.T, request *http.Request) (*http.Response, error)

func MockRoundTripperQueue(t *testing.T, queue []MockRoundTripFunc) http.RoundTripper {
	runNumber := 0

	return MockableRoundTripper{
		RoundTripFunc: func(r *http.Request) (*http.Response, error) {
			defer func() {
				runNumber++
			}()

			if len(queue) <= runNumber {
				t.Fatal("empty queue")
			}

			return queue[runNumber](t, r)
		},
	}
}

type ExpectedTestRequest struct {
	Method string
	Path   string
}

type TestResponse interface {
	CreateResponse() (*http.Response, error)
}

type TestResponseFile struct {
	StatusCode int
	FilePath   string
}

func (f *TestResponseFile) CreateResponse() (*http.Response, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("Response Body File Not Found: %s", f.FilePath)
	}
	defer file.Close()

	return &http.Response{
		StatusCode: f.StatusCode,
		Body:       io.NopCloser(file),
	}, nil
}

func ServeAndValidate(t *testing.T, r TestResponse, expected ExpectedTestRequest) MockRoundTripFunc {
	return func(t *testing.T, request *http.Request) (*http.Response, error) {
		if err := Assert(expected.Method, request.Method); err != nil {
			t.Fatal(err)
		}

		if err := Assert(expected.Path, request.URL.Path); err != nil {
			t.Fatal(err)
		}

		return r.CreateResponse()
	}
}
