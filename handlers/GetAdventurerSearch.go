package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	page := validators.ValidatePageQueryParam(r.URL.Query()["page"])
	queryParams, queryProvided := r.URL.Query()["query"]
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	searchTypeParam := validators.ValidateSearchTypeQueryParam(r.URL.Query()["searchType"])

	if !queryProvided || !validators.ValidateAdventurerName(&queryParams[0]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchType := map[string]uint8{
		"characterName": 1,
		"familyName":    2,
	}[searchTypeParam]

	// Run the scraper
	data, status := scrapers.ScrapeAdventurerSearch(region, queryParams[0], searchType, uint16(page))

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
