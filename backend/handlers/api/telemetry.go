package api

import (
	"errors"
	"net/http"
	"strconv"
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
	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")

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
