package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/tsukanov/beaten-games/data"
)

func main() {
	http.HandleFunc("/", IndexHandler)
	fmt.Println("Starting server on localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
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
