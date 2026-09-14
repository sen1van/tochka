package api

import (
	"errors"
	"net/http"
	"time"
	"tochka/database"
	"tochka/handlers"
	"tochka/middlewares"
)

type newTelemetryRequest struct {
	SensorID  *int       `json:"sensorId" required:"true"`
	TryNumber *int       `json:"tryNumber" required:"true"`
	Value     *int       `json:"value" required:"true"`
	Timestamp *time.Time `json:"timestamp" required:"true"`
}

type telemetryResponse struct {
	Total int                  `json:"total"`
	Data  []database.Telemetry `json:"data"`
}

func postTelemetry(w http.ResponseWriter, r *http.Request) {
	var telemetry newTelemetryRequest

	if !handlers.ReadJSON(w, r, &telemetry) {
		return
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	err = db.NewTelemetry(
		*telemetry.SensorID,
		*telemetry.TryNumber,
		*telemetry.Timestamp,
		*telemetry.Value,
		middlewares.GetAuthToken(r))
	if err != nil {
		if errors.Is(err, database.ErrNoRowAffected) {
			handlers.SendError(w, http.StatusNotFound, err.Error())
		} else {
			handlers.SendError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	handlers.SendJSON(w, http.StatusCreated, nil)
}

func getTelemetry(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := pagination(r)
	if err != nil {
		if errors.Is(err, ErrPagination) {
			handlers.SendError(w, http.StatusBadRequest, err.Error())
		} else {
			handlers.SendError(w, http.StatusInternalServerError, "Pagination issue")
		}

		return
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	pings, total, err := db.GetTelemetry(middlewares.GetAuthToken(r), limit, offset)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	resp := telemetryResponse{
		Total: total,
		Data:  pings,
	}

	handlers.SendJSON(w, http.StatusOK, resp)
}
