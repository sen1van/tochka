package api

import (
	"net/http"
)

func GetApiHandler() http.Handler {
	h := http.NewServeMux()

	h.HandleFunc("GET /ping", getPing)
	h.HandleFunc("POST /telemetry", postTelemetry)

	return h
}
