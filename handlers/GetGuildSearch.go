package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bdo-rest-api/scrapers"
)

func GetGuildSearch(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	regionParams, regionProvided := r.URL.Query()["region"]
	pageParams, pageProvided := r.URL.Query()["page"]
	queryParams, queryProvided := r.URL.Query()["query"]

	// Return status 400 if a required parameter is omitted
	if !queryProvided {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set defaults for optional parameters
	region := defaultRegion
	page := defaultPage

	if regionProvided && validateRegion(regionParams[0]) {
		region = regionParams[0]
	}

	if pageProvided && validatePage(pageParams[0]) {
		page, _ = strconv.Atoi(pageParams[0])
	}

	// Run the scraper
	if data, status := scrapers.ScrapeGuildSearch(region, queryParams[0], int32(page)); status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
