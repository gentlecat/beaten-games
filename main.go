package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/tsukanov/beaten-games/data"
)

func main() {
	fmt.Println("Starting server on localhost:8080...")
	err := http.ListenAndServe(":8080", makeRouter())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func makeRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/games/{id:[0-9]+}", gameHandler)
	r.HandleFunc("/games/add", addHandler).Methods("GET", "POST")
	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	games, err := data.GetAllGames()
	if err != nil {
		log.Fatal("Failed to get games.", err)
	}
	t.Execute(w, struct {
		Games []data.Game
	}{
		games,
	})
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Failed to parse game ID.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// TODO: Implement game lookup
	t, _ := template.ParseFiles("templates/game.html")
	t.Execute(w, struct {
		ID int
	}{
		id,
	})
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/add.html")
		t.Execute(w, nil)
	} else { // POST
		// TODO: Add new game and redirect to its page
	}
}
