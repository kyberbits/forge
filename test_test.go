package forge_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type HandlerTestCase struct {
	Handler            http.Handler
	Request            *http.Request
	ExpectedStatusCode int
	ExpectedBody       string
}

func testHandler(t *testing.T, testCase HandlerTestCase) {
	recorder := httptest.NewRecorder()
	testCase.Handler.ServeHTTP(recorder, testCase.Request)

	actualStatusCode := recorder.Code
	actualBody := recorder.Body.String()

	if actualStatusCode != testCase.ExpectedStatusCode {
		t.Errorf("Got %d, Expected: %d", actualStatusCode, testCase.ExpectedStatusCode)
		return
	}

	if actualBody != testCase.ExpectedBody {
		t.Errorf("Got %s, Expected: %s", actualBody, testCase.ExpectedBody)
		return
	}

}
