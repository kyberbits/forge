package forge

import "net/http"

type HTTPHeaders map[string]string

func (httpHeaders HTTPHeaders) ModifyRequest(r *http.Request) error {
	for k, v := range httpHeaders {
		r.Header.Set(k, v)
	}

	return nil
}

func (httpHeaders HTTPHeaders) CheckResponse(r *http.Response) error {
	return nil
}
