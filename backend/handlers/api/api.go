package api

import (
	"net/http"
)

func GetApiHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /ping", getPing)
	handler.HandleFunc("POST /telemetry", postTelemetry)
	handler.HandleFunc("GET /telemetry", getTelemetry)

	return handler
}
