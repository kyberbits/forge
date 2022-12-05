package forgetest_test

import (
	"net/http"
	"testing"

	"github.com/kyberbits/forge/forgetest"
)

func TestMockRoundTripperQueue(t *testing.T) {
	q := forgetest.MockRoundTripperQueue(nil, []forgetest.MockRoundTripFunc{
		func(t *testing.T, request *http.Request) (*http.Response, error) {
			return nil, nil
		},
	})

	// Request 1
	q.RoundTrip(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	// Request 2 (should panic)
	q.RoundTrip(nil)
}
