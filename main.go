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

	"gitlab.com/man90/black-desert-social-rest-api/scraper"
)

type responseCache struct {
	time time.Time
	data interface{}
}

type errorResponse struct {
	Error string `json:"error"`
}

var (
	globalCacheMap map[string]responseCache = make(map[string]responseCache)
	lastCacheCleanUp time.Time = time.Now()
	cacheTTL time.Duration = time.Hour * 2
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v0/guildProfile", getGuildProfile).Methods("GET")
	router.HandleFunc("/v0/profile", getProfile).Methods("GET")
	router.HandleFunc("/v0/guildProfileSearch", getGuildProfileSearch).Methods("GET")
	router.HandleFunc("/v0/profileSearch", getProfileSearch).Methods("GET")

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

func cleanUpCache() {
	counter := 0
	log.Printf("Cleaning up the cache.")

	for key, element := range globalCacheMap {
		if time.Now().Sub(element.time) > cacheTTL {
			delete(globalCacheMap, key)
			counter++
		}
	}

	log.Printf("%v entries removed from the cache.", counter)
	lastCacheCleanUp = time.Now()
}

func getCachedResponse(cacheMapKey string) (interface{}, bool) {
	cachedReponse, ok := globalCacheMap[cacheMapKey]

	if ok {
		if time.Now().Sub(cachedReponse.time) < cacheTTL {
			log.Printf("Serving \"%v\" from cache.\n", cacheMapKey)
			return cachedReponse.data, true
		} else {
			log.Printf("\"%v\" cache has expired.", cacheMapKey)
		}
	} else {
		log.Printf("\"%v\" not found in cache.", cacheMapKey)
	}
	
	return nil, false
}

func setCachedResponse(cacheMapKey string, data interface{}) {
	if time.Now().Sub(lastCacheCleanUp) > cacheTTL {
		go cleanUpCache()
	}

	globalCacheMap[cacheMapKey] = responseCache{
		time: time.Now(),
		data: data,
	}

	log.Printf("Adding \"%v\" to cache.", cacheMapKey)
}

func getGuildProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	regionParams, ok1 := r.URL.Query()["region"]
	guildNameParams, ok2 := r.URL.Query()["guildName"]

	if !ok1 || !validateRegion(regionParams[0]) || !ok2 {
		return
	}

	cacheMapKey := fmt.Sprintf("getGuildProfile+%v+%v", regionParams[0], guildNameParams[0])
	if cachedReponseData, ok := getCachedResponse(cacheMapKey); ok {
		json.NewEncoder(w).Encode(cachedReponseData)
	} else {
		if data, err := scraper.ScrapeGuildProfile(regionParams[0], guildNameParams[0]); err == nil {
			setCachedResponse(cacheMapKey, data)
			json.NewEncoder(w).Encode(data)
		} else {
			json.NewEncoder(w).Encode(errorResponse{ err.Error() })
		}
	}
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	profileTargetParams, ok := r.URL.Query()["profileTarget"]

	if !ok {
		return
	}

	cacheMapKey := fmt.Sprintf("getProfile+%v", profileTargetParams[0])
	if cachedReponseData, ok := getCachedResponse(cacheMapKey); ok {
		json.NewEncoder(w).Encode(cachedReponseData)
	} else {
		if data, err := scraper.ScrapeProfile(url.QueryEscape(profileTargetParams[0])); err == nil {
			setCachedResponse(cacheMapKey, data)
			json.NewEncoder(w).Encode(data)
		} else {
			json.NewEncoder(w).Encode(errorResponse{ err.Error() })
		}
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

	cacheMapKey := fmt.Sprintf("getGuildProfileSearch+%v+%v+%v", regionParams[0], query, int32(page))
	if cachedReponseData, ok := getCachedResponse(cacheMapKey); ok {
		json.NewEncoder(w).Encode(cachedReponseData)
	} else {
		if data, err := scraper.ScrapeGuildProfileSearch(regionParams[0], query, int32(page)); err == nil {
			setCachedResponse(cacheMapKey, data)
			json.NewEncoder(w).Encode(data)
		} else {
			json.NewEncoder(w).Encode(errorResponse{ err.Error() })
		}
	}
}

func getProfileSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	cacheMapKey := fmt.Sprintf("getProfileSearch+%v+%v+%v+%v", regionParams[0], query, searchType, int32(page))
	if cachedReponseData, ok := getCachedResponse(cacheMapKey); ok {
		json.NewEncoder(w).Encode(cachedReponseData)
	} else {
		if data, err := scraper.ScrapeProfileSearch(regionParams[0], query, searchType, int32(page)); err == nil {
			setCachedResponse(cacheMapKey, data)
			json.NewEncoder(w).Encode(data)
		} else {
			json.NewEncoder(w).Encode(errorResponse{ err.Error() })
		}
	}
}
