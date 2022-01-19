package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
)

func GetGuild(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	regionParams, regionProvided := r.URL.Query()["region"]
	nameParams, nameProvided := r.URL.Query()["guildName"]

	// Return status 400 if a required parameter is omitted
	if !nameProvided {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Set defaults for optional parameters
	region := defaultRegion

	if regionProvided && validateRegion(regionParams[0]) {
		region = regionParams[0]
	}

	// Run the scraper
	if data, status := scrapers.ScrapeGuild(region, nameParams[0]); status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
