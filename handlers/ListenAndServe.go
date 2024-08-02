package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bdo-rest-api/config"
)

func ListenAndServe() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1", getStatus)
	mux.HandleFunc("GET /v1/adventurer", getAdventurer)
	mux.HandleFunc("GET /v1/adventurer/search", getAdventurerSearch)
	mux.HandleFunc("GET /v1/cache", getCache)
	mux.HandleFunc("GET /v1/guild", getGuild)
	mux.HandleFunc("GET /v1/guild/search", getGuildSearch)
	mux.HandleFunc("/", catchall)

	log.Println("Listening for requests")
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", config.GetPort()),
		Handler:      mux,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
