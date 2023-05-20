package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func ReadAllButRetain(reader io.ReadCloser) []byte {
	bodyBytes, _ := io.ReadAll(reader)
	reader.Close()

	// Put the bytes back
	// reader = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

type HTTPClientMiddleware interface {
	ModifyRequest(r *http.Request) error
	CheckResponse(r *http.Response) error
}

type HTTPClientResponseMiddleware interface {
	ModifyResponse(r *http.Response) error
}

func NewHTTPClient(
	httpClient *http.Client,
	middleware []HTTPClientMiddleware,
) *HTTPClient {
	return &HTTPClient{
		httpClient: httpClient,
		middleware: middleware,
	}
}

type HTTPClient struct {
	httpClient *http.Client
	middleware []HTTPClientMiddleware
}

var ErrHTTPClientShouldRetry = errors.New("")

func (c *HTTPClient) Do(request *http.Request) (*http.Response, error) {
	for _, middleware := range c.middleware {
		if err := middleware.ModifyRequest(request); err != nil {
			return nil, err
		}
	}

	const maxAttemptCount = 2

	attemptCount := 0

	var response *http.Response

	for attemptCount < maxAttemptCount {
		attemptCount++

		var err error

		response, err = c.httpClient.Do(request)
		if err != nil {
			return nil, err
		}

		shouldRetry := false

		for _, middleware := range c.middleware {
			if err := middleware.CheckResponse(response); err != nil {
				if errors.Is(err, ErrHTTPClientShouldRetry) {
					shouldRetry = true

					break
				}

				return nil, err
			}
		}

		if !shouldRetry {
			break
		}
	}

	return response, nil
}

func (c *HTTPClient) JSONRequest(
	ctx context.Context,
	method string,
	url string,
	payload any,
	target any,
) error {
	var body io.Reader

	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		body = bytes.NewReader(payloadBytes)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	response, err := c.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if target != nil {
		if err := decoder.Decode(target); err != nil {
			return err
		}
	}

	return nil
}
