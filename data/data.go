package data

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "database.sqlite3"
	sqlInit    = `
		CREATE TABLE IF NOT EXISTS games (
			name TEXT NOT NULL PRIMARY KEY,
			note TEXT,
			beaten_on DATE
		);`
)

type Game struct {
	Name     string
	Note     string
	BeatenOn time.Time
}

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

func AddGame(game Game) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO games (name, note, beaten_on) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(game.Name, game.Note, game.BeatenOn)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func GetAllGames() (games []Game, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, note, beaten_on FROM games ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var curr Game
		err := rows.Scan(&curr.Name, &curr.Note, &curr.BeatenOn)
		if err != nil {
			return nil, err
		}
		games = append(games, curr)
	}
	rows.Close()
	return games, nil
}
