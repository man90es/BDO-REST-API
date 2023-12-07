package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetGuildSearch(w http.ResponseWriter, r *http.Request) {
	page := validators.ValidatePageQueryParam(r.URL.Query()["page"])
	queryParams, queryProvided := r.URL.Query()["query"]
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !queryProvided || !validators.ValidateGuildName(&queryParams[0]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Run the scraper
	data, status := scrapers.ScrapeGuildSearch(region, queryParams[0], uint16(page))

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
