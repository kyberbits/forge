package forge

import "net/http"

type HTTPBasicAuth struct {
	Username string
	Password string
}

func (httpBasicAuth HTTPBasicAuth) ModifyRequest(r *http.Request) error {
	r.SetBasicAuth(httpBasicAuth.Username, httpBasicAuth.Password)

	return nil
}

func (httpBasicAuth HTTPBasicAuth) CheckResponse(r *http.Response) error {
	return nil
}
