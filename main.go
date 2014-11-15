package main

import (
	"fmt"
	"io"
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

	// Regular pages
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/games/{id:[0-9]+}", gameHandler)
	r.HandleFunc("/games/add", addHandler).Methods("GET", "POST")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	return r
}

// executeTemplates is a custom tempate executor that uses our template
// structure. Should be used when rendering templates based on "base.html"
// template.
func executeTemplates(wr io.Writer, data interface{}, filenames ...string) error {
	filenames = append(filenames, "templates/base.html")
	t, err := template.ParseFiles(filenames...)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(wr, "base", data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	games, err := data.GetAllGames()
	if err != nil {
		http.Error(w, "Failed to get games.", http.StatusInternalServerError)
		return
	}
	err = executeTemplates(w, struct{ Games []data.Game }{games},
		"templates/index.html")
	if err != nil {
		http.Error(w, "Failed to execute template.", http.StatusInternalServerError)
		return
	}
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Failed to parse game ID.", http.StatusInternalServerError)
		return
	}

	// TODO: Implement game lookup

	err = executeTemplates(w, struct{ ID int }{id}, "templates/game.html")
	if err != nil {
		http.Error(w, "Failed to execute template.", http.StatusInternalServerError)
		return
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := executeTemplates(w, nil, "templates/add.html")
		if err != nil {
			http.Error(w, "Failed to execute template.", http.StatusInternalServerError)
			return
		}

	} else { // POST (new game is sumbitted)
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse submitted form.", http.StatusInternalServerError)
			return
		}
		vals := r.Form
		var game data.Game
		game.Name = vals.Get("name")
		game.Note = vals.Get("note")
		//game.BeatenOn = r.Form["beaten_on"]

		err = data.AddGame(game)
		if err != nil {
			http.Error(w, "Failed to add a game.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
