package api

import (
	"net/http"
	"time"
	"tochka/handlers"
	"tochka/middlewares"
)

// ```
// POST /api/telemetry
// BODY:
//
//	json: {
//	  "sensor_id": 12345,
//	  "try_number": 12345,
//	  "value": 12345,
//	  "timestamp": "2025-01-01T00:00:00Z"
//	}
//
// ANSW:
//
//	201 Created
//
// ```

type telemetryReq struct {
	SensorID  *int       `json:"sensor_id" required:"true"`
	TryNumber *int       `json:"try_number" required:"true"`
	Value     *int       `json:"value" required:"true"`
	Timestamp *time.Time `json:"timestamp" required:"true"`
}

func postTelemetry(w http.ResponseWriter, r *http.Request) {
	var telemetry telemetryReq

	if !handlers.ReadJSON(w, r, &telemetry) {
		return
	}

	db, err := middlewares.GetDB(r)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = db.NewTelemetry(*telemetry.SensorID, *telemetry.TryNumber, *telemetry.Timestamp)
	if err != nil {
		handlers.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handlers.SendJSON(w, http.StatusCreated, nil)
}
