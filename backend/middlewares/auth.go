package middlewares

import "net/http"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := GetAuthToken(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetAuthToken(r *http.Request) string {
	return r.Header.Get("X-API-Token")
}
