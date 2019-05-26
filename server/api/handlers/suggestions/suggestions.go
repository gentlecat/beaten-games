package suggestions

import (
	"encoding/json"
	bomb "go.roman.zone/go-bomb"
	"log"
	"net/http"
)

func SuggestGamesHandler(gbClient *bomb.GBClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}
