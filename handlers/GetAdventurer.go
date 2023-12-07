package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetAdventurer(w http.ResponseWriter, r *http.Request) {
	profileTarget, profileTargetGood := validators.ValidateProfileTargetQueryParam(r.URL.Query()["profileTarget"])
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !profileTargetGood {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Run the scraper
	data, status := scrapers.ScrapeAdventurer(region, profileTarget)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
