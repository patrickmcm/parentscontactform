package middleware

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
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
