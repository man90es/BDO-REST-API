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

	if data, status, date, expires, ok := cache.ProfileSearch.GetRecord([]string{region, query, searchType}); ok {
		w.Header().Set("Expires", expires)
		w.Header().Set("Last-Modified", date)

		if status == http.StatusOK {
			json.NewEncoder(w).Encode(data)
		} else {
			w.WriteHeader(status)
		}

		return
	}

	if ok := giveMaintenanceResponse(w, region); ok {
		return
	}

	if tasksQuantityExceeded := scraper.EnqueueAdventurerSearch(r.Header.Get("CF-Connecting-IP"), region, query, searchType); tasksQuantityExceeded {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "You have exceeded the maximum number of concurrent tasks.",
		})

		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Player search is being fetched. Please try again later.",
	})
}
