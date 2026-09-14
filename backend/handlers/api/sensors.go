package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"tochka/database"
	"tochka/handlers"
	"tochka/middlewares"
)

const maxLimit = 100

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
		if errors.Is(err, database.ErrNoRowAffected) {
			handlers.SendError(w, http.StatusConflict, "Sensor already exists")

			return
		}

		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	handlers.SendJSON(w, http.StatusCreated, nil)
}

func getSensors(w http.ResponseWriter, r *http.Request) {
	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	sensors, err := db.GetSensors(middlewares.GetAuthToken(r))
	if err != nil {
		slog.Error("failed to get sensors", "error", err)
		handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")

		return
	}

	handlers.SendJSON(w, http.StatusOK, sensorsResponse{
		Total: len(sensors),
		Data:  sensors,
	})
}

func getSensorTelemetry(w http.ResponseWriter, r *http.Request) {
	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")
	sensorIDQuery := r.PathValue("id")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			handlers.SendError(w, http.StatusBadRequest, "limit must be a number")

			return
		}

		limit = 20
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			handlers.SendError(w, http.StatusBadRequest, "offset must be a number")

			return
		}

		offset = 0
	}

	if limit > maxLimit {
		handlers.SendError(w, http.StatusBadRequest, "Limit must be less than 100")

		return
	}

	if limit == 0 {
		handlers.SendError(w, http.StatusBadRequest, "Limit must be greater than 0")

		return
	}

	if limit < 0 || offset < 0 {
		handlers.SendError(w, http.StatusBadRequest, "Limit and offset must be non-negative")

		return
	}

	sensorID, err := strconv.Atoi(sensorIDQuery)
	if err != nil {
		handlers.SendError(w, http.StatusBadRequest, "Invalid sensor ID")

		return
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

	err = db.UpdateSensor(middlewares.GetAuthToken(r), sensorID, *req.Name)
	if err != nil {
		if errors.Is(err, database.ErrNoRowAffected) {
			handlers.SendError(w, http.StatusNotFound, "Sensor not found")
		} else {
			handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")
		}

		return
	}

	handlers.SendJSON(w, http.StatusNoContent, nil)
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

	err = db.DeleteSensor(middlewares.GetAuthToken(r), sensorID)
	if err != nil {
		if errors.Is(err, database.ErrNoRowAffected) {
			handlers.SendError(w, http.StatusNotFound, "Sensor not found")
		} else {
			handlers.SendError(w, http.StatusInternalServerError, "Some troubles with db")
		}

		return
	}

	handlers.SendJSON(w, http.StatusNoContent, nil)
}
