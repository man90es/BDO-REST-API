package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetAdventurer(w http.ResponseWriter, r *http.Request) {
	profileTargetParams, profileTargetProvided := r.URL.Query()["profileTarget"]
	regionParams, regionProvided := r.URL.Query()["region"]

	// Return status 400 if a required parameter is invalid
	if !profileTargetProvided || !validators.ValidateProfileTarget(&profileTargetParams[0]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set defaults for optional parameters
	region := defaultRegion

	if regionProvided && validators.ValidateRegion(&regionParams[0]) {
		region = regionParams[0]
	}

	// Run the scraper
	data, status := scrapers.ScrapeAdventurer(region, profileTargetParams[0])

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
