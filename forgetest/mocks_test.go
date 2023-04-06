package forgetest_test

import (
	"net/http"
	"testing"

	"github.com/kyberbits/forge/forgetest"
)

func TestMockRoundTripperQueue(t *testing.T) {
	q := forgetest.MockRoundTripperQueue(nil, []forgetest.MockRoundTripFunc{
		func(t *testing.T, request *http.Request) (*http.Response, error) {
			response := &http.Response{
				Body: http.NoBody,
			}

			return response, nil
		},
	})

	// Request 1
	response1, _ := q.RoundTrip(nil)
	defer response1.Body.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	// Request 2 (should panic)
	response2, _ := q.RoundTrip(nil)
	defer response2.Body.Close()
}
