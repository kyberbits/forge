package forge_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/kyberbits/forge/forge"
	"github.com/kyberbits/forge/forgetest"
)

func TestHTTPAuthOAuth2(t *testing.T) {
	ctx := context.Background()
	httpClient := &http.Client{
		Transport: forgetest.MockRoundTripperQueue(t, []forgetest.MockRoundTripFunc{
			forgetest.ServeAndValidate(
				t,
				&forgetest.TestResponseFile{
					StatusCode: http.StatusOK,
					FilePath:   "test_files/client/oauth2_200_access_token.json",
				},
				forgetest.ExpectedTestRequest{
					Method: http.MethodPost,
					Path:   "/auth",
				},
			),
			forgetest.ServeAndValidate(
				t,
				&forgetest.TestResponseFile{
					StatusCode: http.StatusOK,
					FilePath:   "test_files/client/test.json",
				},
				forgetest.ExpectedTestRequest{
					Method: http.MethodGet,
					Path:   "/api/test",
				},
			),
			forgetest.ServeAndValidate(
				t,
				&forgetest.TestResponseFile{
					StatusCode: http.StatusOK,
					FilePath:   "test_files/client/test.json",
				},
				forgetest.ExpectedTestRequest{
					Method: http.MethodGet,
					Path:   "/api/test",
				},
			),
			forgetest.ServeAndValidate(
				t,
				&forgetest.TestResponseFile{
					StatusCode: http.StatusUnauthorized,
					FilePath:   "test_files/client/test.json",
				},
				forgetest.ExpectedTestRequest{
					Method: http.MethodGet,
					Path:   "/api/test",
				},
			),
			forgetest.ServeAndValidate(
				t,
				&forgetest.TestResponseFile{
					StatusCode: http.StatusOK,
					FilePath:   "test_files/client/oauth2_200_access_token.json",
				},
				forgetest.ExpectedTestRequest{
					Method: http.MethodPost,
					Path:   "/auth",
				},
			),
			forgetest.ServeAndValidate(
				t,
				&forgetest.TestResponseFile{
					StatusCode: http.StatusOK,
					FilePath:   "test_files/client/test.json",
				},
				forgetest.ExpectedTestRequest{
					Method: http.MethodGet,
					Path:   "/api/test",
				},
			),
		}),
	}

	oauth2 := forge.NewHTTPOAuth2(httpClient, forge.HTTPOAuth2GrantPassword{
		ClientID:     "fake-client-id",
		ClientSecret: "fake-client-secret",
		Username:     "username",
		Password:     "password",
	}, "/auth")

	c := forge.NewHTTPClient(
		httpClient,
		[]forge.HTTPClientMiddleware{
			oauth2,
		},
	)

	if err := c.JSONRequest(
		ctx,
		http.MethodGet,
		"/api/test",
		nil,
		nil,
	); err != nil {
		t.Error(err)
	}

	if err := c.JSONRequest(
		ctx,
		http.MethodGet,
		"/api/test",
		nil,
		nil,
	); err != nil {
		t.Error(err)
	}

	if err := c.JSONRequest(
		ctx,
		http.MethodGet,
		"/api/test",
		nil,
		nil,
	); err != nil {
		t.Error(err)
	}
}
