package forge_test

import (
	"net/http"
	"testing"

	"github.com/kyberbits/forge/forge"
)

func TestRouter(t *testing.T) {
	type Response struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}

	router := &forge.HTTPRouter{
		Routes: map[string]http.Handler{
			"/": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerContext := forge.NewHandlerContext(w, r)
				handlerContext.RespondJSON(http.StatusOK, Response{
					OK:      true,
					Message: "Hello there.",
				})
			}),
		},
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerContext := forge.NewHandlerContext(w, r)

			handlerContext.RespondJSON(http.StatusNotFound, Response{
				OK:      false,
				Message: "not found",
			})
		}),
	}

	{ // Check root
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			panic(err)
		}
		testHandler(t, HandlerTestCase{
			Handler:            router,
			Request:            request,
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       `{"ok":true,"message":"Hello there."}` + "\n",
		})
	}

	{ // Check 404
		request, err := http.NewRequest(http.MethodGet, "/foobar", nil)
		if err != nil {
			panic(err)
		}
		testHandler(t, HandlerTestCase{
			Handler:            router,
			Request:            request,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedBody:       `{"ok":false,"message":"not found"}` + "\n",
		})
	}

	{ // Check custom 404
		request, err := http.NewRequest(http.MethodGet, "/foobar", nil)
		if err != nil {
			panic(err)
		}
		testHandler(t, HandlerTestCase{
			Handler:            router,
			Request:            request,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedBody:       `{"ok":false,"message":"not found"}` + "\n",
		})
	}
}
