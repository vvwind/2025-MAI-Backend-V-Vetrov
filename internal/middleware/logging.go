package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the request
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Create a response writer wrapper to capture status code
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		// Log the response
		latency := time.Since(start)
		log.Printf(
			"Completed %s %s | Status: %d | Duration: %v",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			latency,
		)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
