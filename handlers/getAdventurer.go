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

func getAdventurer(w http.ResponseWriter, r *http.Request) {
	profileTarget, profileTargetOk, profileTargetValidationMessage := validators.ValidateProfileTargetQueryParam(r.URL.Query()["profileTarget"])
	if !profileTargetOk {
		giveBadRequestResponse(w, profileTargetValidationMessage)
		return
	}

	region, regionOk, regionValidationMessage := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	if !regionOk {
		giveBadRequestResponse(w, regionValidationMessage)
		return
	}

	if ok := giveMaintenanceResponse(w, region); ok {
		return
	}

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := cache.Profiles.GetRecord([]string{region, profileTarget})
	if !found {
		go scrapers.ScrapeAdventurer(region, profileTarget)

		var wg sync.WaitGroup
		wg.Add(1)

		// TODO: Maintenance handling if it was detected while waiting for the scraper
		cache.Profiles.Bus.Subscribe(strings.Join([]string{region, profileTarget}, ","), func(v cache.CacheEntry[models.Profile]) {
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
