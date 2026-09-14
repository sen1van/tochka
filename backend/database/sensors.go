package database

import (
	"fmt"
	"log/slog"
)

type Sensor struct {
	ID   int
	Name string
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

func (db *DB) ChangeSensorName(id int, name string) error {
	ctx, cancel := timeoutContext()
	defer cancel()

	_, err := db.db.ExecContext(ctx, `
		UPDATE sensors
			SET name = ?
			WHERE id = ?
	`, name, id)
	if err != nil {
		return fmt.Errorf("failed to change sensor name: %w", err)
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

func (db *DB) GetSensors(token string) ([]Sensor, error) {
	ctx, cancel := timeoutContext()
	defer cancel()

	rows, err := db.db.QueryContext(ctx, `
		SELECT id, name
		FROM sensors
		WHERE owner_token = ?
	`, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensors: %w", err)
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

		err := rows.Scan(&sensor.ID, &sensor.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sensor: %w", err)
		}

		sensors = append(sensors, sensor)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate sensors: %w", err)
	}

	return sensors, nil
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
