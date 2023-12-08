package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	page := validators.ValidatePageQueryParam(r.URL.Query()["page"])
	query, queryOk := validators.ValidateAdventurerNameQueryParam(r.URL.Query()["query"])
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	searchTypeParam := validators.ValidateSearchTypeQueryParam(r.URL.Query()["searchType"])

	if !queryOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchType := map[string]uint8{
		"characterName": 1,
		"familyName":    2,
	}[searchTypeParam]

	data, status := scrapers.ScrapeAdventurerSearch(region, query, searchType, page)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
