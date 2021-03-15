package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"
	"strconv"
	"net/url"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"

	"gitlab.com/man90/black-desert-social-rest-api/scraper"
)

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1000000),
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(2 * time.Hour),
		cache.ClientWithRefreshKey("opn"),
	)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	router := mux.NewRouter()

	router.Handle("/v0/guildProfile", 		cacheClient.Middleware(http.HandlerFunc(getGuildProfile)))		.Methods("GET")
	router.Handle("/v0/profile", 			cacheClient.Middleware(http.HandlerFunc(getProfile)))			.Methods("GET")
	router.Handle("/v0/guildProfileSearch", cacheClient.Middleware(http.HandlerFunc(getGuildProfileSearch))).Methods("GET")
	router.Handle("/v0/profileSearch", 		cacheClient.Middleware(http.HandlerFunc(getProfileSearch)))		.Methods("GET")

	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8001"
	}

	srv := &http.Server{
		Handler: 		router,
		Addr: 			fmt.Sprintf("0.0.0.0:%v", port),
		WriteTimeout: 	15 * time.Second,
		ReadTimeout:  	15 * time.Second,
	}

	log.Printf("Listening for requests on port %v.", port)
	log.Fatal(srv.ListenAndServe())
}

func validateRegion(r string) bool {
	return r == "EU" || r == "NA"
}

func getGuildProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	regionParams, ok1 := r.URL.Query()["region"]
	guildNameParams, ok2 := r.URL.Query()["guildName"]

	if !ok1 || !validateRegion(regionParams[0]) || !ok2 {
		return
	}

	if data, err := scraper.ScrapeGuildProfile(regionParams[0], guildNameParams[0]); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(errorResponse{ err.Error() })
	}
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	profileTargetParams, ok := r.URL.Query()["profileTarget"]

	if !ok {
		return
	}

	if data, err := scraper.ScrapeProfile(url.QueryEscape(profileTargetParams[0])); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(errorResponse{ err.Error() })
	}
}

func getGuildProfileSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	if data, err := scraper.ScrapeGuildProfileSearch(regionParams[0], query, int32(page)); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(errorResponse{ err.Error() })
	}
}

func getProfileSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	regionParams, ok1 := r.URL.Query()["region"]
	searchTypeParams, ok2 := r.URL.Query()["searchType"]
	pageParams, ok3 := r.URL.Query()["page"]
	queryParams, ok4 := r.URL.Query()["query"]

	searchType := map[string]int8 {
		"characterName": 1,
		"familyName": 2,
	}[searchTypeParams[0]]

	if !ok1 || !ok2 || !ok4 || !validateRegion(regionParams[0]) || searchType < 1 || searchType > 2 {
		return
	}

	page := 1

	if ok3 {
		page, _ = strconv.Atoi(pageParams[0])
	}

	if data, err := scraper.ScrapeProfileSearch(regionParams[0], queryParams[0], searchType, int32(page)); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(errorResponse{ err.Error() })
	}
}
