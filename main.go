package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dlintw/goconf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.roman.zone/beaten-games/data"
	bomb "go.roman.zone/go-bomb"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

var (
	gbClient *bomb.GBClient
)

func main() {
	fmt.Println("Loading configuration...")
	config, err := goconf.ReadConfigFile("config.txt")
	if err != nil {
		log.Fatal("Failed to load config file! ", err)
	}
	apiKey, err := config.GetString("default", "giant_bomb_api_key")
	if err != nil {
		log.Fatal("Failed to get Giant Bomb API key from config file!", err)
	}

	gbClient = bomb.NewClient(apiKey)

	router := makeRouter()
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	srv := &http.Server{
		Handler: loggedRouter,
		Addr:    "localhost:8080",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on http://localhost:8080...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func makeRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", indexHandler)

	r.HandleFunc("/games/{id:[0-9]+}", gameHandler)
	r.HandleFunc("/api/games/add", addHandler).Methods("GET", "POST")
	r.HandleFunc("/api/games/quick-add", quickAddHandler).Methods("POST")
	r.HandleFunc("/api/games", getGamesHandler).Methods("GET")
	r.HandleFunc("/api/games/delete", deleteHandler).Methods("POST")

	r.HandleFunc("/api/suggest/games", suggestGamesHandler)

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("frontend/static"))))

	return r
}

// executeTemplates is a custom template executor that uses our template
// structure. Should be used when rendering templates based on "base.html"
// template.
func executeTemplates(wr io.Writer, data interface{}, filenames ...string) error {
	filenames = append(filenames, "frontend/templates/base.html")
	t, err := template.ParseFiles(filenames...)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(wr, "base", data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := executeTemplates(w, struct{}{}, "frontend/templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse submitted form.", http.StatusInternalServerError)
		return
	}
	vals := r.Form
	var game data.Game
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

	err = data.AddGame(game)
	if err != nil {
		http.Error(w, "Failed to add a game.", http.StatusInternalServerError)
		return
	}

}

func quickAddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse submitted form.", http.StatusInternalServerError)
		return
	}
	vals := r.Form
	var game data.Game
	game.Name = vals.Get("name")
	game.Note = sql.NullString{
		Valid: false,
	}
	game.BeatenOn = data.NullTime{
		Valid: false,
	}

	err = data.AddGame(game)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add a game.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse submitted form.", http.StatusInternalServerError)
		return
	}
	vals := r.Form
	rowsAffected, err := data.DeleteGame(vals.Get("name"))
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

func suggestGamesHandler(w http.ResponseWriter, r *http.Request) {
	qvals := r.URL.Query()
	query, ok := qvals["q"]
	if !ok {
		http.Error(w, "Query is empty.", http.StatusBadRequest)
		return
	}

	resp, err := gbClient.Search(query[0], 10, 1, []string{bomb.ResourceGame}, nil)
	if err != nil {
		http.Error(w, "Search failed.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	b, err := json.Marshal(resp.Results)
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := data.GetAllGames()

	b, err := json.Marshal(games)
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
