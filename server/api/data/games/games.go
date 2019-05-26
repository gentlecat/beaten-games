package games

import (
	"database/sql"
	data2 "go.roman.zone/beaten-games/server/api/data"
)

type GameEntity struct {
	Name     string
	Note     sql.NullString
	BeatenOn data2.NullTime
}

func AddGame(game GameEntity) error {
	db, err := data2.OpenDB()
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
	db, err := data2.OpenDB()
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

func GetAllGames() (games []GameEntity, err error) {
	db, err := data2.OpenDB()
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
		var curr GameEntity
		err := rows.Scan(&curr.Name, &curr.Note, &curr.BeatenOn)
		if err != nil {
			return nil, err
		}
		games = append(games, curr)
	}
	return games, nil
}
