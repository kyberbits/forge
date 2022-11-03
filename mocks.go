package forge

import (
	"net/http"
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
				panic("empty queue")
			}

			return queue[runNumber](t, r)
		},
	}
}
