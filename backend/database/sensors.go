package database

import (
	"fmt"
	"log/slog"
	"time"
)

type Sensor struct {
	ID            int        `json:"sensorId"`
	Name          string     `json:"name"`
	LastTelemetry *Telemetry `json:"lastTelemetry"`
}

func (db *DB) createSensorTable() error {
	ctx, cancel := timeoutContext()
	defer cancel()

	_, err := db.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS sensors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			owner_token TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create sensor table: %w", err)
	}

	return nil
}

func (db *DB) NewSensor(name, token string, id int) error {
	ctx, cancel := timeoutContext()
	defer cancel()

	_, err := db.db.ExecContext(ctx, `
		INSERT INTO sensors (id, name, owner_token)
			VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET name = ?, owner_token = ?;
	`, id, name, token, name, token)
	if err != nil {
		return fmt.Errorf("failed to create sensor owner: %w", err)
	}

	return nil
}

func (db *DB) DeleteSensor(id int) error {
	ctx, cancel := timeoutContext()
	defer cancel()

	_, err := db.db.ExecContext(ctx, `
		DELETE FROM sensors
			WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete sensor: %w", err)
	}

	return nil
}

func (db *DB) GetSensorsCount(token string) (int, error) {
	ctx, cancel := timeoutContext()
	defer cancel()

	var count int
	err := db.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
			FROM sensors
			WHERE owner_token = ?
	`, token).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get sensor count: %w", err)
	}

	return count, nil
}

func (db *DB) GetSensors(token string, limit int, offset int) ([]Sensor, int, error) {
	ctx, cancel := timeoutContext()
	defer cancel()

	rows, err := db.db.QueryContext(ctx, `
	WITH ranked_telemetry AS (
	    SELECT
	        sensor_id,
	        try,
	        timestamp,
	        value,
	        ROW_NUMBER() OVER (PARTITION BY sensor_id ORDER BY timestamp DESC) as rn
	    FROM telemetry
		)
	SELECT
	    s.id,
	    s.name,
	    t.timestamp,
	    t.try,
	    t.value
			FROM sensors s
	LEFT JOIN ranked_telemetry t ON t.sensor_id = s.id AND t.rn = 1
		WHERE s.owner_token = ?
		LIMIT ? OFFSET ?
	`, token, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get sensors: %w", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("closing rows", "error", err)
		}
	}()

	var sensors []Sensor

	for rows.Next() {
		var sensor Sensor

		var (
			telemetryTime  *time.Time
			telemetryTry   *int
			telemetryValue *int
		)

		err := rows.Scan(
			&sensor.ID, &sensor.Name,
			&telemetryTime, &telemetryTry, &telemetryValue,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sensor: %w", err)
		}

		if telemetryValue != nil {
			sensor.LastTelemetry = &Telemetry{
				SensorID:  sensor.ID,
				Timestamp: *telemetryTime,
				Try:       *telemetryTry,
				Value:     *telemetryValue,
			}
		}
		sensors = append(sensors, sensor)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to iterate sensors: %w", err)
	}

	sensorCount, err := db.GetSensorsCount(token)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get sensor count: %w", err)
	}
	return sensors, sensorCount, nil
}

func (db *DB) UpdateSensor(id int, name string) error {
	ctx, cancel := timeoutContext()
	defer cancel()

	_, err := db.db.ExecContext(ctx, `
		UPDATE sensors
			SET name = ?
			WHERE id = ?
	`, name, id)
	if err != nil {
		return fmt.Errorf("failed to update sensor: %w", err)
	}

	return nil
}
