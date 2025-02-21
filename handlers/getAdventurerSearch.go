package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/scraper"
	"bdo-rest-api/validators"
)

func getAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	region, regionOk, regionValidationMessage := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	if !regionOk {
		giveBadRequestResponse(w, regionValidationMessage)
		return
	}

	query, queryOk, queryValidationMessage := validators.ValidateAdventurerNameQueryParam(r.URL.Query()["query"], region)
	if !queryOk {
		giveBadRequestResponse(w, queryValidationMessage)
		return
	}

	searchTypeQueryParam := r.URL.Query()["searchType"]
	searchType := validators.ValidateSearchTypeQueryParam(searchTypeQueryParam)

	if ok := giveMaintenanceResponse(w, region); ok {
		return
	}

	data, status, date, expires, found := cache.ProfileSearch.GetRecord([]string{region, query, searchType})
	if !found {
		go scraper.EnqueueAdventurerSearch(region, query, searchType)
		data, status, date, expires = cache.ProfileSearch.WaitForRecord([]string{region, query, searchType})
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
