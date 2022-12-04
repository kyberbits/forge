package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HTTPClientInterceptor interface {
	ModRequest(request *http.Request) error
	ErrorCheck(response *http.Response) error
}

func NewHTTPClient(
	interceptor HTTPClientInterceptor,
	client *http.Client,
) *HTTPClient {
	return &HTTPClient{
		client:      client,
		interceptor: interceptor,
	}
}

type HTTPClient struct {
	client      *http.Client
	interceptor HTTPClientInterceptor
}

func (httpClient *HTTPClient) Do(request *http.Request) (*http.Response, error) {
	// Run the ModRequest func
	if err := httpClient.interceptor.ModRequest(request); err != nil {
		return nil, err
	}

	response, err := httpClient.client.Do(request)
	if err != nil {
		return nil, err
	}

	// Copy the body so we can re-write it to the response
	bodyBytes, _ := io.ReadAll(response.Body)
	response.Body.Close()

	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Run the ErrorCheck func
	if err := httpClient.interceptor.ErrorCheck(response); err != nil {
		return nil, err
	}

	// Replace the body incase the error checker read the body
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return response, nil
}

func (httpClient *HTTPClient) Request(
	ctx context.Context,
	method string,
	url string,
	payload interface{},
	target interface{},
) error {
	// Build the body, if there is one
	var body io.Reader
	if payload != nil {
		payloadBytes, _ := json.Marshal(payload)
		body = bytes.NewReader(payloadBytes)
	}

	// Build the http.Request
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	request = request.WithContext(ctx)

	// Execute the http.Request
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	// Decode the response body onto the target, if there is one
	decoder := json.NewDecoder(response.Body)
	if target != nil {
		if err := decoder.Decode(target); err != nil {
			return err
		}
	}

	return nil
}
