package database

func (db *DB) createDeviceTable() error {
	_, err := db.db.Exec(`
		CREATE TABLE IF NOT EXISTS device (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			owner_token TEXT,
		);
	`)
	return err
}

func (db *DB) NewOwnerDevice(name, token string) error {
	_, err := db.db.Exec(`
		INSERT INTO device (name, owner_token)
			VALUES (?, ?)
		ON CONFLICT(id) DO UPDATE SET name = ?, owner_token = ?;
	`, name, token, name, token)
	return err
}
