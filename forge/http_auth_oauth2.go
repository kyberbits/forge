package forge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type HTTPOAuth2Token struct {
	AccessToken string `json:"access_token"`
}

type HTTPOAuth2Grant interface {
	GenerateQuery() (url.Values, error)
}

func NewHTTPOAuth2(
	client *http.Client,
	grant HTTPOAuth2Grant,
	tokenURL string,
) *HTTPOAuth2 {
	return &HTTPOAuth2{
		client:   client,
		grant:    grant,
		mutex:    &sync.Mutex{},
		token:    nil,
		tokenURL: tokenURL,
	}
}

type HTTPOAuth2 struct {
	client   *http.Client
	grant    HTTPOAuth2Grant
	mutex    *sync.Mutex
	token    *HTTPOAuth2Token
	tokenURL string
}

func (oauth2 *HTTPOAuth2) ModifyRequest(r *http.Request) error {
	// Make sure we have a valid token
	if oauth2.token == nil {
		// Lock the mutex to make sure every request does not get their own token
		oauth2.mutex.Lock()
		defer oauth2.mutex.Unlock()

		// Check again after we get past the lock incase another routine set the token for us
		if oauth2.token == nil {
			newToken, err := oauth2.fetchNewToken(r.Context())
			if err != nil {
				return err
			}

			oauth2.token = newToken
		}
	}

	// Add the access token to the request
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauth2.token.AccessToken))

	return nil
}

func (oauth2 *HTTPOAuth2) fetchNewToken(ctx context.Context) (*HTTPOAuth2Token, error) {
	tokenRequestValues, err := oauth2.grant.GenerateQuery()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s?%s", oauth2.tokenURL, tokenRequestValues.Encode()),
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := oauth2.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("bad status code: %d", response.StatusCode)
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	newToken := &HTTPOAuth2Token{}
	if err := decoder.Decode(newToken); err != nil {
		return nil, err
	}

	return newToken, nil
}

func (oauth2 *HTTPOAuth2) CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusUnauthorized {
		oauth2.token = nil

		return ErrHTTPClientShouldRetry
	}

	return nil
}

type HTTPOAuth2GrantPassword struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	Scope        string
}

func (grant HTTPOAuth2GrantPassword) GenerateQuery() (url.Values, error) {
	values := url.Values{}

	values.Set("grant_type", "password")
	values.Set("client_id", grant.ClientID)
	values.Set("client_secret", grant.ClientSecret)
	values.Set("username", grant.Username)
	values.Set("password", grant.Password)

	return values, nil
}
