package middlewares

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf(
			"method=%s path=%s status=%d duration=%v",
			r.Method,
			r.URL.Path,
			rw.status,
			time.Since(start),
		)
	})
}
