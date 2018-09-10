package log

import (
	"context"
	"net/http"
	"time"
)

// StatusResponseWriter implements the http.ResponseWriter interface, by wrapping another
// http.ResponseWriter. It has the ability to record the HTTP response code for logging
// purposes.
type StatusResponseWriter struct {
	status int
	http.ResponseWriter
}

// Status returns the ResponseWriter's HTTP status code
func (srw *StatusResponseWriter) Status() int {
	return srw.status
}

// WriteHeader writes a HTTP status code to the client
func (srw *StatusResponseWriter) WriteHeader(status int) {
	srw.status = status
	srw.ResponseWriter.WriteHeader(status)
}

// LoggingHandler logs incoming HTTP requests and their responses
func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextKeyRequestMethod, r.Method)
		ctx = context.WithValue(ctx, ContextKeyEndpoint, r.URL.Path)
		ctx = context.WithValue(ctx, ContextKeyRequestHost, r.Host)
		Info(ctx, "Inbound request")
		statusWriter := StatusResponseWriter{200, w}
		r = r.WithContext(ctx)
		startTime := time.Now()
		next.ServeHTTP(&statusWriter, r)
		endTime := time.Now()
		Info(ctx, "Outbound response", Duration("duration", endTime.Sub(startTime)), Int("status", statusWriter.Status()))
	})
}
