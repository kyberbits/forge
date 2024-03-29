package forge_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/kyberbits/forge/forge"
)

func TestStatic(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/foo/bar.txt", nil)
	if err != nil {
		panic(err)
	}

	testHandler(t, HandlerTestCase{
		Handler: &forge.HTTPStatic{
			FileSystem: http.FS(os.DirFS("test_files/static")),
		},
		Request:            request,
		ExpectedStatusCode: http.StatusOK,
		ExpectedBody:       "foobar\n",
	})
}

func TestStaticIndex(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		panic(err)
	}

	testHandler(t, HandlerTestCase{
		Handler: &forge.HTTPStatic{
			FileSystem: http.FS(os.DirFS("test_files/static")),
		},
		Request:            request,
		ExpectedStatusCode: http.StatusOK,
		ExpectedBody:       "Hello there.\n",
	})
}

func TestStaticNotFound(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/not-found", nil)
	if err != nil {
		panic(err)
	}

	testHandler(t, HandlerTestCase{
		Handler: &forge.HTTPStatic{
			FileSystem: http.FS(os.DirFS("test_files/static")),
		},
		Request:            request,
		ExpectedStatusCode: http.StatusNotFound,
		ExpectedBody:       "404 page not found\n",
	})
}

func TestStaticCustomNotFound(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/not-found", nil)
	if err != nil {
		panic(err)
	}

	testHandler(t, HandlerTestCase{
		Handler: &forge.HTTPStatic{
			FileSystem: http.FS(os.DirFS("test_files/static")),
		},
		Request:            request,
		ExpectedStatusCode: http.StatusNotFound,
		ExpectedBody:       "404 page not found\n",
	})
}
