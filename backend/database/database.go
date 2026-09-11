package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func NewDB(link string) *DB {
	db, err := sql.Open("sqlite", link)
	if err != nil {
		panic(err)
	}
	d := &DB{db: db}

	d.createDeviceTable()
	d.createTelemetryTable()

	return d
}

func (d *DB) Close() error {
	return d.db.Close()
}
