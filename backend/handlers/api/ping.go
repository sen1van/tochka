package api

import (
	"net/http"
	"tochka/handlers"
)

func getPing(w http.ResponseWriter, r *http.Request) {
	handlers.SendJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}
