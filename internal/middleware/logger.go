package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture the status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler in the chain
		next.ServeHTTP(wrapped, r)

		// Log after the request is handled
		duration := time.Since(start)
		if wrapped.statusCode < 400 || wrapped.statusCode == 404 {
			log.Printf(
				"%s %s %d %s",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
			)
		}
	}
}

func LogAndError(r *http.Request, w http.ResponseWriter, friendlyErrMsg string, errMsg string, code int) {
	log.Printf(
		"ERROR: %s %s %d %s %s %s %s %s %s",
		r.Method,
		r.URL.Path,
		code,
		r.RemoteAddr,
		r.UserAgent(),
		r.Referer(),
		friendlyErrMsg,
		errMsg,
	)
	http.Error(w, friendlyErrMsg, code)
}
