package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	regionParams, regionProvided := r.URL.Query()["region"]
	searchTypeParams, searchTypeProvided := r.URL.Query()["searchType"]
	pageParams, pageProvided := r.URL.Query()["page"]
	queryParams, queryProvided := r.URL.Query()["query"]

	// Return status 400 if a required parameter is invalid
	if !queryProvided || !validators.ValidateAdventurerName(&queryParams[0]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set defaults for optional parameters
	region := defaultRegion
	searchType := uint8(2)
	page := defaultPage

	if regionProvided && validators.ValidateRegion(&regionParams[0]) {
		region = regionParams[0]
	}

	if searchTypeProvided && validators.ValidateSearchType(&searchTypeParams[0]) {
		searchType = map[string]uint8{
			"characterName": 1,
			"familyName":    2,
		}[searchTypeParams[0]]
	}

	if pageProvided && validators.ValidatePage(&pageParams[0]) {
		page, _ = strconv.Atoi(pageParams[0])
	}

	// Run the scraper
	data, status := scrapers.ScrapeAdventurerSearch(region, queryParams[0], searchType, uint16(page))

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
