package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func NewHTTPClient(
	client *http.Client,
	modRequest func(request *http.Request) error,
	errorCheck func(response *http.Response) error,
) *HTTPClient {
	return &HTTPClient{
		client:     client,
		modRequest: modRequest,
		errorCheck: errorCheck,
	}
}

type HTTPClient struct {
	client     *http.Client
	modRequest func(request *http.Request) error
	errorCheck func(response *http.Response) error
}

func (httpClient *HTTPClient) Do(request *http.Request) (*http.Response, error) {
	// Run the ModRequest func
	if httpClient.modRequest != nil {
		if err := httpClient.modRequest(request); err != nil {
			return nil, err
		}
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
	if httpClient.errorCheck != nil {
		if err := httpClient.errorCheck(response); err != nil {
			return nil, err
		}
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
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

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
	defer response.Body.Close()

	// Decode the response body onto the target, if there is one
	decoder := json.NewDecoder(response.Body)
	if target != nil {
		if err := decoder.Decode(target); err != nil {
			return err
		}
	}

	return nil
}
