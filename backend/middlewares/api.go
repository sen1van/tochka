package middlewares

import (
	"context"
	"errors"
	"net/http"
	"tochka/database"
)

type contextKey string

const (
	DBContextKey contextKey = "db"
)

var ErrDBNotFound = errors.New("database not found in context")

func APIMiddleware(next http.Handler, db *database.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx := context.WithValue(r.Context(), DBContextKey, db)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetDB(r *http.Request) (*database.DB, error) {
	value, ok := r.Context().Value(DBContextKey).(*database.DB)
	if !ok {
		return nil, ErrDBNotFound
	}

	return value, nil
}
