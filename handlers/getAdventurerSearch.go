package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"bdo-rest-api/cache"
	"bdo-rest-api/scrapers"
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

	// All names are non-case-sensitive, so this will allow to utilise cache better
	query = strings.ToLower(query)

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := cache.ProfileSearch.GetRecord([]string{region, query, searchType})
	if !found {
		data, status = scrapers.ScrapeAdventurerSearch(region, query, searchType)

		if status == http.StatusInternalServerError {
			w.WriteHeader(status)
			return
		}

		if ok := giveMaintenanceResponse(w, region); ok {
			return
		}

		date, expires = cache.ProfileSearch.AddRecord([]string{region, query, searchType}, data, status)
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
