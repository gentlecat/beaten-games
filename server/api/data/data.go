package data

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "database.sqlite3"
	sqlInit    = `
		CREATE TABLE IF NOT EXISTS games (
			name TEXT NOT NULL PRIMARY KEY,
			note TEXT,
			beaten_on TIMESTAMP
			added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
)

// OpenDB opens database and, if successful, returns a reference to it.
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	// TODO: Check if database file exists before attempting to create a table.
	_, err = db.Exec(sqlInit)
	if err != nil {
		return nil, err
	}
	return db, nil
}
