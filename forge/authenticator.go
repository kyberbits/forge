package forge

import (
	"net/http"
)

type User interface {
	Username() string
	Roles() []string
}

type Authenticator interface {
	LoginHandler() http.Handler
	InitializeSession(next http.Handler) http.Handler
	RequireLogin(next http.Handler) http.Handler
	GetUser(r *http.Request) (User, error)
}
