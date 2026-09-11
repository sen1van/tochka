package database

import "fmt"

func (db *DB) createDeviceTable() error {
	_, err := db.db.ExecContext(timeoutContext(), `
		CREATE TABLE IF NOT EXISTS device (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			owner_token TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create device table: %w", err)
	}

	return nil
}

func (db *DB) NewOwnerDevice(name, token string) error {
	_, err := db.db.ExecContext(timeoutContext(), `
		INSERT INTO device (name, owner_token)
			VALUES (?, ?)
		ON CONFLICT(id) DO UPDATE SET name = ?, owner_token = ?;
	`, name, token, name, token)
	if err != nil {
		return fmt.Errorf("failed to create owner device: %w", err)
	}

	return nil
}
