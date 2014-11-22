package data

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "database.sqlite3"
	sqlInit    = `
		CREATE TABLE IF NOT EXISTS games (
			name TEXT NOT NULL PRIMARY KEY,
			note TEXT,
			beaten_on TIMESTAMP
		);`
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

type Game struct {
	Name     string
	Note     sql.NullString
	BeatenOn NullTime
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

	if game.BeatenOn.Valid { // No time specified
		stmt, err := tx.Prepare("INSERT INTO games (name, note, beaten_on) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(game.Name, game.Note, game.BeatenOn.Time)
		if err != nil {
			return err
		}
		tx.Commit()
	} else {
		stmt, err := tx.Prepare("INSERT INTO games (name, note) VALUES (?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(game.Name, game.Note)
		if err != nil {
			return err
		}
		tx.Commit()
	}

	return nil
}

func DeleteGame(name string) (int64, error) {
	log.Println(name)
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare("DELETE FROM games WHERE name = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name)
	if err != nil {
		return 0, err
	}
	tx.Commit()
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func GetAllGames() (games []Game, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, note, beaten_on FROM games ORDER BY beaten_on DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var curr Game
		err := rows.Scan(&curr.Name, &curr.Note, &curr.BeatenOn.Time)
		if err != nil {
			return nil, err
		}
		games = append(games, curr)
	}
	rows.Close()
	return games, nil
}
