package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"

	"github.com/gorilla/mux"

	"gitlab.com/man90/black-desert-social-rest-api/scraper"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v0/guildProfile", getGuildProfile).Methods("GET")
	// router.HandleFunc("/v0/profile", getProfile).Methods("GET")

	srv := &http.Server{
		Handler: 		router,
		Addr: 			"127.0.0.1:8001",
		WriteTimeout: 	15 * time.Second,
		ReadTimeout:  	15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func getGuildProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	regionParams, ok1 := r.URL.Query()["region"]
	guildNameParams, ok2 := r.URL.Query()["guildName"]

	if !ok1 || !ok2 {
		return
	}

	json.NewEncoder(w).Encode(scraper.ScrapeGuildProfile(regionParams[0], guildNameParams[0]))
}

// func getProfile(w http.ResponseWriter, r *http.Request) {}
