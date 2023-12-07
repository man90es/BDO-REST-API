package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetGuild(w http.ResponseWriter, r *http.Request) {
	nameParams, nameProvided := r.URL.Query()["guildName"]
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !nameProvided || !validators.ValidateGuildName(&nameParams[0]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Run the scraper
	data, status := scrapers.ScrapeGuild(region, nameParams[0])

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
