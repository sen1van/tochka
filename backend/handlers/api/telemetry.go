package api

import (
	"net/http"
	"strconv"
	"time"
	"tochka/database"
	"tochka/handlers"
	"tochka/middlewares"
)

type newTelemetryRequest struct {
	SensorID  *int       `json:"sensor_id" required:"true"`
	TryNumber *int       `json:"try_number" required:"true"`
	Value     *int       `json:"value" required:"true"`
	Timestamp *time.Time `json:"timestamp" required:"true"`
}

type telemetryResp struct {
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

	err = db.NewTelemetry(*telemetry.SensorID, *telemetry.TryNumber, *telemetry.Timestamp, middlewares.GetAuthToken(r))
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	handlers.SendJSON(w, http.StatusCreated, nil)
}

func getTelemetry(w http.ResponseWriter, r *http.Request) {
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
		handlers.SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	pings, total, err := db.GetTelemetry(middlewares.GetAuthToken(r), limit, offset)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, err.Error())

		return
	}

	resp := telemetryResp{
		Total: total,
		Data:  pings,
	}

	handlers.SendJSON(w, http.StatusOK, resp)
}
