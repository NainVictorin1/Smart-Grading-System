package main

import (
	"log"
	"net/http"
	"time"
)

// loggingMiddleware logs the details of each incoming HTTP request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Create a custom response writer to capture the status code
		ww := &responseWriter{ResponseWriter: w}

		// Call the next handler
		next.ServeHTTP(ww, r)

		// Log the status code and the time it took to process the request
		log.Printf("Completed %s %s with status %d in %v", r.Method, r.URL.Path, ww.statusCode, time.Since(start))
	})
}

// responseWriter is a custom HTTP response writer that captures the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code of the response.
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
