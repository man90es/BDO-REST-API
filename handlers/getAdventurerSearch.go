package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/utils"
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

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := cache.ProfileSearch.GetRecord([]string{region, query, searchType})
	if !found {
		go scrapers.ScrapeAdventurerSearch(region, query, searchType)

		var wg sync.WaitGroup
		wg.Add(1)

		// TODO: Maintenance handling if it was detected while waiting for the scraper
		cache.ProfileSearch.Bus.Subscribe(strings.Join([]string{region, query, searchType}, ","), func(v cache.CacheEntry[[]models.Profile]) {
			data = v.Data
			date = utils.FormatDateForHeaders(v.Date)
			expires = utils.FormatDateForHeaders(v.Expires)
			status = v.Status

			wg.Done()
		})

		wg.Wait()
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
