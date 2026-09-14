package database

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

const (
	defaultTimeout = 5 * time.Second
)

type DB struct {
	db *sql.DB
}

func NewDB(link string) *DB {
	db, err := sql.Open("sqlite", link)
	if err != nil {
		panic(err)
	}

	newDB := &DB{db: db}

	err = newDB.createSensorTable()
	if err != nil {
		panic(err)
	}

	err = newDB.createTelemetryTable()
	if err != nil {
		panic(err)
	}

	return newDB
}

func (d *DB) Close() {
	slog.Error("closing database", "error", d.db.Close())
}

func timeoutContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	return ctx, cancel
}
