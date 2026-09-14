package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"tochka/database"
	"tochka/handlers"
	"tochka/middlewares"
)

type newSensorRequest struct {
	SensorID *int    `json:"sensorId" required:"true"`
	Name     *string `json:"name" required:"true"`
}

type updateSensorRequest struct {
	Name *string `json:"name" required:"true"`
}

type sensorsResponse struct {
	Total int               `json:"total"`
	Data  []database.Sensor `json:"data"`
}

func postSensor(w http.ResponseWriter, r *http.Request) {
	var req newSensorRequest

	if !handlers.ReadJSON(w, r, &req) {
		return
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	err = db.NewSensor(*req.Name, middlewares.GetAuthToken(r), *req.SensorID)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	handlers.SendJSON(w, http.StatusCreated, nil)
}

func getSensors(w http.ResponseWriter, r *http.Request) {
	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		offset = 0
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	sensors, total, err := db.GetSensors(middlewares.GetAuthToken(r), limit, offset)
	if err != nil {
		slog.Error("failed to get sensors", "error", err)
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	handlers.SendJSON(w, http.StatusOK, sensorsResponse{
		Total: total,
		Data:  sensors,
	})
}

func getSensorTelemetry(w http.ResponseWriter, r *http.Request) {
	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")
	sensorIDQuery := r.PathValue("id")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		offset = 0
	}

	sensorID, err := strconv.Atoi(sensorIDQuery)
	if err != nil {
		sensorID = 0
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	telemetry, total, err := db.GetSensorTelemetry(middlewares.GetAuthToken(r), sensorID, limit, offset)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	resp := telemetryResponse{
		Total: total,
		Data:  telemetry,
	}

	handlers.SendJSON(w, http.StatusOK, resp)
}

func updateSensor(w http.ResponseWriter, r *http.Request) {
	idReq := r.PathValue("id")

	sensorID, err := strconv.Atoi(idReq)
	if err != nil {
		handlers.SendError(w, http.StatusBadRequest, "Invalid sensor ID")

		return
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	var req updateSensorRequest
	if !handlers.ReadJSON(w, r, &req) {
		handlers.SendError(w, http.StatusBadRequest, "Invalid request body")

		return
	}

	if req.Name == nil {
		handlers.SendError(w, http.StatusBadRequest, "Name is required")

		return
	}

	err = db.UpdateSensor(sensorID, *req.Name)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	handlers.SendJSON(w, http.StatusOK, nil)
}

func deleteSensor(w http.ResponseWriter, r *http.Request) {
	idReq := r.PathValue("id")

	sensorID, err := strconv.Atoi(idReq)
	if err != nil {
		handlers.SendError(w, http.StatusBadRequest, "Invalid sensor ID")

		return
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	err = db.DeleteSensor(sensorID)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	handlers.SendJSON(w, http.StatusOK, nil)
}
