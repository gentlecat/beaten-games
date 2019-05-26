package game

import (
	"encoding/json"
	gamesData "go.roman.zone/beaten-games/server/api/data/games"
	"log"
	"net/http"
	"time"
)

type GameResponse struct {
	Name     string    `json:"name"`
	Note     string    `json:"note"`
	BeatenOn time.Time `json:"beaten_on"`
}

func GetGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := gamesData.GetAllGames()
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var output []GameResponse

	for _, game := range games {
		resp := GameResponse{
			Name: game.Name,
		}
		if game.Note.Valid {
			resp.Note = game.Note.String
		}
		if game.BeatenOn.Valid {
			resp.BeatenOn = game.BeatenOn.Time
		}

		output = append(output, resp)
	}

	b, err := json.Marshal(output)
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
