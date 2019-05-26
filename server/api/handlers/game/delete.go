package game

import (
	gamesData "go.roman.zone/beaten-games/server/api/data/games"
	"log"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse submitted form.", http.StatusInternalServerError)
		return
	}
	vals := r.Form
	rowsAffected, err := gamesData.DeleteGame(vals.Get("name"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete a game.", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Can't find this game.", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
