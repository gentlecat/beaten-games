package api

import (
	"fmt"
	"github.com/dlintw/goconf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	gameHandlers "go.roman.zone/beaten-games/server/api/handlers/game"
	suggestionsHandlers "go.roman.zone/beaten-games/server/api/handlers/suggestions"
	bomb "go.roman.zone/go-bomb"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func CreateServer(config *goconf.ConfigFile) *http.Server {
	router := makeRouter(config)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	return &http.Server{
		Handler: loggedRouter,
		Addr:    "localhost:8080",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func makeRouter(config *goconf.ConfigFile) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", indexHandler)

	r.HandleFunc("/games/{id:[0-9]+}", gameHandler)
	r.HandleFunc("/api/games/add", gameHandlers.AddHandler).Methods("GET", "POST")
	r.HandleFunc("/api/games/quick-add", gameHandlers.QuickAddHandler).Methods("POST")
	r.HandleFunc("/api/games", gameHandlers.GetGamesHandler).Methods("GET")
	r.HandleFunc("/api/games/delete", gameHandlers.DeleteHandler).Methods("POST")

	apiKey, err := config.GetString("default", "giant_bomb_api_key")
	if err != nil {
		log.Fatal("Failed to get Giant Bomb API key from config file!", err)
	}
	r.Handle("/api/suggest/games", suggestionsHandlers.SuggestGamesHandler(bomb.NewClient(apiKey)))

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
