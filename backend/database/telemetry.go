package database

import (
	"fmt"
	"log/slog"
	"time"
)

type Telemetry struct {
	ID        int       `json:"id"`
	DeviceID  int       `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	Try       int       `json:"try"`
	Token     string    `json:"token"`
}

func (db *DB) createTelemetryTable() error {
	_, err := db.db.ExecContext(timeoutContext(), `
		CREATE TABLE IF NOT EXISTS telemetry (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER,
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

func (db *DB) NewTelemetry(deviceID int, try int, timestamp time.Time, token string) error {
	query := `INSERT INTO telemetry (device_id, try, timestamp, token) VALUES (?, ?, ?, ?)`

	_, err := db.db.ExecContext(timeoutContext(), query, deviceID, try, timestamp.UTC(), token)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}

func (db *DB) getTelemetryCount(token string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM telemetry
		WHERE device_id IN (SELECT device_id FROM device WHERE owner_token = ?)`

	rows, err := db.db.QueryContext(timeoutContext(), query, token)
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
		SELECT id, device_id, timestamp, try, token
		FROM telemetry
		WHERE device_id IN (SELECT device_id FROM device WHERE owner_token = ?)
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?`

	rows, err := db.db.QueryContext(timeoutContext(), query, token, limit, offset)
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

		err := rows.Scan(&ping.ID, &ping.DeviceID, &ping.Timestamp, &ping.Try, &ping.Token)
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
