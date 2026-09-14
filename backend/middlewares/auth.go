package middlewares

import (
	"net/http"
	"tochka/handlers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := GetAuthToken(r)
		if token == "" {
			handlers.SendError(w, http.StatusUnauthorized, "Unauthorized")

			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetAuthToken(r *http.Request) string {
	return r.Header.Get("X-Api-Token")
}
