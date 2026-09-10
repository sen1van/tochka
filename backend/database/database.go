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
	d.Connect(link)

	d.createDeviceTable()
	d.createPingTable()

	return d
}

func (d *DB) Connect(link string) error {
	var err error
	d.db, err = sql.Open("sqlite", link)
	return err
}

func (d *DB) Close() error {
	return d.db.Close()
}
