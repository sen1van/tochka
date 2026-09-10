package database

type Ping struct {
	deviceID  int
	timestamp string
	try       int
}

func (db *DB) createPingTable() error {
	_, err := db.db.Exec(`
		CREATE TABLE IF NOT EXISTS ping (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			try INTEGER,
		);
	`)
	return err
}

func (db *DB) NewPing(deviceID int, try int) error {
	query := `INSERT INTO ping (device_id, try) VALUES (?, ?)`
	_, err := db.db.Exec(query, deviceID, try)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetPings(token string) ([]Ping, error) {
	query := `SELECT * FROM ping WHERE device_id in (SELECT device_id FROM device WHERE owner_token = ?)`
	rows, err := db.db.Query(query, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pings []Ping

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		var p Ping
		if err := rows.Scan(&p.deviceID, &p.timestamp, &p.try); err != nil {
			return nil, err
		}
		pings = append(pings, p)
	}
	return pings, nil
}
