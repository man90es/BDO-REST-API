package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetGuild(w http.ResponseWriter, r *http.Request) {
	regionParams, regionProvided := r.URL.Query()["region"]
	nameParams, nameProvided := r.URL.Query()["guildName"]

	// Return status 400 if a required parameter is invalid
	if !nameProvided || !validators.ValidateGuildName(&nameParams[0]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set defaults for optional parameters
	region := defaultRegion

	if regionProvided && validators.ValidateRegion(&regionParams[0]) {
		region = regionParams[0]
	}

	// Run the scraper
	data, status := scrapers.ScrapeGuild(region, nameParams[0])

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
