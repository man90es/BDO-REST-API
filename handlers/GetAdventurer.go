package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"bdo-rest-api/scrapers"
)

func GetAdventurer(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	profileTargetParams, profileTargetProvided := r.URL.Query()["profileTarget"]

	// Return status 400 if a required parameter is omitted
	if !profileTargetProvided {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Run the scraper
	if data, status := scrapers.ScrapeAdventurer(url.QueryEscape(profileTargetParams[0])); status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
