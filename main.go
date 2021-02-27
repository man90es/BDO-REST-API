package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"
	"strconv"
	"net/url"

	"github.com/gorilla/mux"

	"gitlab.com/man90/black-desert-social-rest-api/scraper"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v0/guildProfile", getGuildProfile).Methods("GET")
	router.HandleFunc("/v0/profile", getProfile).Methods("GET")
	router.HandleFunc("/v0/guildProfileSearch", getGuildProfileSearch).Methods("GET")
	router.HandleFunc("/v0/profileSearch", getProfileSearch).Methods("GET")

	srv := &http.Server{
		Handler: 		router,
		Addr: 			"127.0.0.1:8001",
		WriteTimeout: 	15 * time.Second,
		ReadTimeout:  	15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func validateRegion(r string) bool {
	return r == "EU" || r == "NA"
}

func getGuildProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	regionParams, ok1 := r.URL.Query()["region"]
	guildNameParams, ok2 := r.URL.Query()["guildName"]

	if !ok1 || !validateRegion(regionParams[0]) || !ok2 {
		return
	}

	json.NewEncoder(w).Encode(scraper.ScrapeGuildProfile(regionParams[0], guildNameParams[0]))
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	profileTargetParams, ok := r.URL.Query()["profileTarget"]

	if !ok {
		return
	}

	json.NewEncoder(w).Encode(scraper.ScrapeProfile(url.QueryEscape(profileTargetParams[0])))
}

func getGuildProfileSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	regionParams, ok1 := r.URL.Query()["region"]
	pageParams, ok2 := r.URL.Query()["page"]
	queryParams, ok3 := r.URL.Query()["query"]

	if !ok1 || !validateRegion(regionParams[0]) {
		return
	}

	page := 1

	if ok2 {
		page, _ = strconv.Atoi(pageParams[0])
	}

	var query string

	if ok3 {
		query = queryParams[0]
	}

	json.NewEncoder(w).Encode(scraper.ScrapeGuildProfileSearch(regionParams[0], query, int32(page)))
}

func getProfileSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	regionParams, ok1 := r.URL.Query()["region"]
	searchTypeParams, ok2 := r.URL.Query()["searchType"]
	pageParams, ok3 := r.URL.Query()["page"]
	queryParams, ok4 := r.URL.Query()["query"]

	if !ok1 || !validateRegion(regionParams[0]) {
		return
	}

	var searchType int8

	if ok2 {
		sT := map[string]int8 {
			"characterName": 1,
			"familyName": 2,
		}

		var ok bool

		if searchType, ok = sT[searchTypeParams[0]]; !ok {
			searchType = 3
		}
	}

	page := 1

	if ok3 {
		page, _ = strconv.Atoi(pageParams[0])
	}

	var query string

	if ok4 {
		query = queryParams[0]
	}

	json.NewEncoder(w).Encode(scraper.ScrapeProfileSearch(regionParams[0], query, searchType, int32(page)))
}
