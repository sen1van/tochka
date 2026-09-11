package database

import (
	"time"
)

type Ping struct {
	ID        int       `json:"id"`
	DeviceID  int       `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	Try       int       `json:"try"`
}

func (db *DB) createTelemetryTable() error {
	_, err := db.db.Exec(`
		CREATE TABLE IF NOT EXISTS telemetry (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER,
			timestamp DATETIME,
			try INTEGER
		);
	`)
	return err
}

func (db *DB) NewTelemetry(deviceID int, try int, timestamp time.Time) error {
	query := `INSERT INTO telemetry (device_id, try, timestamp) VALUES (?, ?, ?)`
	_, err := db.db.Exec(query, deviceID, try, timestamp.UTC())
	return err
}

func (db *DB) GetTelemetry(token string) ([]Ping, error) {
	query := `
		SELECT id, device_id, timestamp, try
		FROM telemetry
		WHERE device_id IN (SELECT device_id FROM device WHERE owner_token = ?)`

	rows, err := db.db.Query(query, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pings []Ping

	for rows.Next() {
		var p Ping
		if err := rows.Scan(&p.ID, &p.DeviceID, &p.Timestamp, &p.Try); err != nil {
			return nil, err
		}
		pings = append(pings, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pings, nil
}
