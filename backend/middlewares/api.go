package middlewares

import (
	"context"
	"net/http"
	"tochka/database"
)

func APIMiddleware(next http.Handler, db *database.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(r.Context(), "db", db)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
