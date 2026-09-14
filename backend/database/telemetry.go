package database

import (
	"fmt"
	"log/slog"
	"time"
)

type Telemetry struct {
	ID        int       `json:"id"`
	SensorID  int       `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Try       int       `json:"try"`
	Token     string    `json:"token"`
}

func (db *DB) createTelemetryTable() error {
	ctx, cancel := timeoutContext()
	defer cancel()

	_, err := db.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS telemetry (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sensor_id INTEGER,
			timestamp DATETIME,
			try INTEGER,
			token TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("create table error: %w", err)
	}

	return nil
}

func (db *DB) NewTelemetry(sensorID int, try int, timestamp time.Time, token string) error {
	query := `INSERT INTO telemetry (sensor_id, try, timestamp, token) VALUES (?, ?, ?, ?)`

	ctx, cancel := timeoutContext()
	defer cancel()

	_, err := db.db.ExecContext(ctx, query, sensorID, try, timestamp.UTC(), token)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}

func (db *DB) getTelemetryCount(token string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM telemetry
		WHERE sensor_id IN (SELECT id FROM sensors WHERE owner_token = ?)`

	ctx, cancel := timeoutContext()

	defer cancel()

	rows, err := db.db.QueryContext(ctx, query, token)
	if err != nil {
		return 0, fmt.Errorf("query error: %w", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("closing rows", "error", err)
		}
	}()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scan error: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("rows error: %w", err)
	}

	return count, nil
}

func (db *DB) GetTelemetry(token string, limit int, offset int) ([]Telemetry, int, error) {
	query := `
		SELECT id, sensor_id, timestamp, try, token
		FROM telemetry
		WHERE sensor_id IN (SELECT id FROM sensors WHERE owner_token = ?)
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?`

	ctx, cancel := timeoutContext()
	defer cancel()

	rows, err := db.db.QueryContext(ctx, query, token, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("closing rows", "error", err)
		}
	}()

	var telemetry []Telemetry

	for rows.Next() {
		var ping Telemetry

		err := rows.Scan(&ping.ID, &ping.SensorID, &ping.Timestamp, &ping.Try, &ping.Token)
		if err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		telemetry = append(telemetry, ping)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	count, err := db.getTelemetryCount(token)
	if err != nil {
		return nil, 0, err
	}

	return telemetry, count, nil
}

func (db *DB) getSensorTelemetryCount(token string, sensorID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM telemetry
		WHERE sensor_id IN (SELECT id FROM sensors WHERE owner_token = ? AND id = ?)`

	ctx, cancel := timeoutContext()

	defer cancel()

	rows, err := db.db.QueryContext(ctx, query, token, sensorID)
	if err != nil {
		return 0, fmt.Errorf("query error: %w", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("closing rows", "error", err)
		}
	}()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scan error: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("rows error: %w", err)
	}

	return count, nil
}

func (db *DB) GetSensorTelemetry(token string, sensorID int, limit int, offset int) ([]Telemetry, int, error) {
	query := `
		SELECT id, sensor_id, timestamp, try, token
		FROM telemetry
		WHERE sensor_id IN (SELECT id FROM sensors WHERE owner_token = ?) AND sensor_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?`

	ctx, cancel := timeoutContext()
	defer cancel()

	rows, err := db.db.QueryContext(ctx, query, token, sensorID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("closing rows", "error", err)
		}
	}()

	var telemetry []Telemetry

	for rows.Next() {
		var ping Telemetry

		err := rows.Scan(&ping.ID, &ping.SensorID, &ping.Timestamp, &ping.Try, &ping.Token)
		if err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}

		telemetry = append(telemetry, ping)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	count, err := db.getSensorTelemetryCount(token, sensorID)
	if err != nil {
		return nil, 0, err
	}

	return telemetry, count, nil
}
