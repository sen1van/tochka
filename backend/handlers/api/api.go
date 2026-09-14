package api

import (
	"net/http"
)

func GetApiHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /ping", getPing)

	handler.HandleFunc("POST /telemetry", postTelemetry)
	handler.HandleFunc("GET /telemetry", getTelemetry)

	handler.HandleFunc("GET /sensors/{id}/telemetry", getSensorTelemetry)
	handler.HandleFunc("POST /sensors", postSensor)
	handler.HandleFunc("GET /sensors", getSensors)
	handler.HandleFunc("PATCH /sensors/{id}", updateSensor)
	handler.HandleFunc("DELETE /sensors/{id}", deleteSensor)

	return handler
}
