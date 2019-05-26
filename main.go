package main

import (
	"fmt"
	"github.com/dlintw/goconf"
	"go.roman.zone/beaten-games/server/api"
	"log"
)

func main() {
	fmt.Println("Loading configuration...")
	config, err := goconf.ReadConfigFile("config.txt")
	if err != nil {
		log.Fatal("Failed to load config file! ", err)
	}

	fmt.Println("Starting server on http://localhost:8080...")
	err = api.CreateServer(config).ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
