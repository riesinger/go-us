package log

import (
	"context"
	"net/http"
)

// InboundLoggingHandler logs incoming HTTP requests
func InboundLoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextKeyRequestMethod, r.Method)
		ctx = context.WithValue(ctx, ContextKeyEndpoint, r.URL.Path)
		ctx = context.WithValue(ctx, ContextKeyRequestHost, r.Host)
		Info(ctx, "Inbound request")
	})
}
