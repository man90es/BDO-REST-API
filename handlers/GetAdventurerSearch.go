package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bdo-rest-api/scrapers"
)

func GetAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	regionParams, regionProvided := r.URL.Query()["region"]
	searchTypeParams, searchTypeProvided := r.URL.Query()["searchType"]
	pageParams, pageProvided := r.URL.Query()["page"]
	queryParams, queryProvided := r.URL.Query()["query"]

	// Return status 400 if a required parameter is omitted
	if !queryProvided {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set defaults for optional parameters
	region := defaultRegion
	searchType := uint8(2)
	page := defaultPage

	if regionProvided && validateRegion(regionParams[0]) {
		region = regionParams[0]
	}

	if searchTypeProvided && validateSearchType(searchTypeParams[0]) {
		searchType = map[string]uint8{
			"characterName": 1,
			"familyName":    2,
		}[searchTypeParams[0]]
	}

	if pageProvided && validatePage(pageParams[0]) {
		page, _ = strconv.Atoi(pageParams[0])
	}

	// Run the scraper
	if data, status := scrapers.ScrapeAdventurerSearch(region, queryParams[0], searchType, uint16(page)); status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
