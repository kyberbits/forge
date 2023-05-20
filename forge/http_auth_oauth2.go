package forge

import "net/http"

type HTTPOAuth2Grant interface{}

func NewHTTPOAuth2(
	httpClient *http.Client,
	grant HTTPOAuth2Grant,
) *HTTPOAuth2 {
	return &HTTPOAuth2{
		grant: grant,
	}
}

type HTTPOAuth2 struct {
	grant HTTPOAuth2Grant
}

func (oauth2 *HTTPOAuth2) ModifyRequest(r *http.Request) error {
	return nil
}

func (oauth2 *HTTPOAuth2) CheckResponse(r *http.Response) error {
	return nil
}

type HTTPOAuth2GrantPassword struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	Scope        string
}
