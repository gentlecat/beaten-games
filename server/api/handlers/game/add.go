package game

import (
	"database/sql"
	"go.roman.zone/beaten-games/server/api/data"
	gamesData "go.roman.zone/beaten-games/server/api/data/games"
	"log"
	"net/http"
	"time"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse submitted form.", http.StatusInternalServerError)
		return
	}
	vals := r.Form
	var game gamesData.GameEntity
	game.Name = vals.Get("name")
	game.Note = sql.NullString{
		String: vals.Get("note"),
		Valid:  true,
	}
	if len(r.Form["beaten_on"]) > 0 && len(r.Form["beaten_on"][0]) > 0 {
		parsed, err := time.Parse("2006-01-02", r.Form["beaten_on"][0])
		if err != nil {
			http.Error(w, "Failed to parse date.", http.StatusBadRequest)
			return
		}
		game.BeatenOn = data.NullTime{
			Time:  parsed,
			Valid: true,
		}
	} else {
		game.BeatenOn = data.NullTime{
			Valid: true,
		}
	}

	err = gamesData.AddGame(game)
	if err != nil {
		http.Error(w, "Failed to add a game.", http.StatusInternalServerError)
		return
	}
}

func QuickAddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse submitted form.", http.StatusInternalServerError)
		return
	}
	vals := r.Form
	var game gamesData.GameEntity
	game.Name = vals.Get("name")
	game.Note = sql.NullString{
		Valid: false,
	}
	game.BeatenOn = data.NullTime{
		Valid: false,
	}

	err = gamesData.AddGame(game)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add a game.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
