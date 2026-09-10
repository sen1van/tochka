package handlers

import (
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("OK"))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
