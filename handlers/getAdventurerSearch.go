package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

var profileSearchCache = cache.NewCache[[]models.Profile]()

func getAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	region, regionOk := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	if !regionOk {
		giveBadRequestResponse(w)
		return
	}

	query, queryOk := validators.ValidateAdventurerNameQueryParam(r.URL.Query()["query"], region)
	if !queryOk {
		giveBadRequestResponse(w)
		return
	}

	page := validators.ValidatePageQueryParam(r.URL.Query()["page"])
	searchTypeQueryParam := r.URL.Query()["searchType"]
	searchType, searchTypeAsString := validators.ValidateSearchTypeQueryParam(searchTypeQueryParam)

	if ok := giveMaintenanceResponse(w, region); ok {
		return
	}

	// All names are non-case-sensitive, so this will allow to utilise cache better
	query = strings.ToLower(query)

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := profileSearchCache.GetRecord([]string{region, query, searchTypeAsString, fmt.Sprint(page)})
	if !found {
		data, status = scrapers.ScrapeAdventurerSearch(region, query, searchType, page)

		if status == http.StatusInternalServerError {
			w.WriteHeader(status)
			return
		}

		if ok := giveMaintenanceResponse(w, region); ok {
			return
		}

		date, expires = profileSearchCache.AddRecord([]string{region, query, searchTypeAsString, fmt.Sprint(page)}, data, status)
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
