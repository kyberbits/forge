package forgeutils

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKeyType string

const (
	contextIDKey          contextKeyType = "forge-context-id"
	contextHTTPRequestKey contextKeyType = "forge-request"
)

func ContextAddID(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, contextIDKey, uuid.New().String())

	return ctx
}

func ContextGetID(ctx context.Context) string {
	rawValue := ctx.Value(contextIDKey)

	contextID, ok := rawValue.(string)
	if !ok {
		return ""
	}

	return contextID
}

func ContextGetRequest(ctx context.Context) *http.Request {
	rawValue := ctx.Value(contextHTTPRequestKey)

	httpRequest, ok := rawValue.(*http.Request)
	if !ok {
		return nil
	}

	return httpRequest
}

func ContextAddToRequest(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = ContextAddID(ctx)
	ctx = context.WithValue(ctx, contextHTTPRequestKey, r) // TODO: this seems wrong

	return r.WithContext(ctx)
}
