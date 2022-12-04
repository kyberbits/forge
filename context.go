package forge

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKeyType string

const contextIDKey contextKeyType = "forge-context-id"
const contextHTTPRequestKey contextKeyType = "forge-request"

func AddContextID(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, contextIDKey, uuid.New().String())
	return ctx
}

func GetContextID(ctx context.Context) string {
	rawValue := ctx.Value(contextIDKey)
	contextID, ok := rawValue.(string)
	if !ok {
		return ""
	}

	return contextID
}

func getContextRequest(ctx context.Context) *http.Request {
	rawValue := ctx.Value(contextHTTPRequestKey)
	httpRequest, ok := rawValue.(*http.Request)
	if !ok {
		return nil
	}

	return httpRequest
}

func addContextValuesToRequest(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = AddContextID(ctx)
	ctx = context.WithValue(ctx, contextHTTPRequestKey, r) // TODO: this seems wrong
	return r.WithContext(ctx)
}
