package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	ErrPagination = errors.New("pagination error")
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

func pagination(r *http.Request) (int, int, error) {
	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			return 0, 0, fmt.Errorf("%w: limit must be a number", ErrPagination)
		}

		limit = 20
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			return 0, 0, fmt.Errorf("%w: offset must be a number", ErrPagination)
		}

		offset = 0
	}

	if limit > maxLimit {
		return 0, 0, fmt.Errorf("%w: limit must be less than 100", ErrPagination)
	}

	if limit == 0 {
		return 0, 0, fmt.Errorf("%w: limit must be greater than 0", ErrPagination)
	}

	if limit < 0 || offset < 0 {
		return 0, 0, fmt.Errorf("%w: limit and offset must be non-negative", ErrPagination)
	}

	return limit, offset, nil
}
