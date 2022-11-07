package forge_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/kyberbits/forge"
)

func TestHTTPSpa(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/foo/bar.txt", nil)

	handler := forge.HTTPSpaHandler(
		http.FS(os.DirFS("test_files/static")),
		"index.html",
		func(h http.Header) {},
	)

	testHandler(t, HandlerTestCase{
		Handler:            handler,
		Request:            request,
		ExpectedStatusCode: http.StatusOK,
		ExpectedBody:       "foobar\n",
	})

	request2, _ := http.NewRequest(http.MethodGet, "/fake/page", nil)

	testHandler(t, HandlerTestCase{
		Handler:            handler,
		Request:            request2,
		ExpectedStatusCode: http.StatusOK,
		ExpectedBody:       "Hello there.\n",
	})
}
