package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

var profilesCache = cache.NewCache[models.Profile]()

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
	data, status, date, expires, found := profilesCache.GetRecord([]string{region, profileTarget})
	if !found {
		data, status = scrapers.ScrapeAdventurer(region, profileTarget)

		if status == http.StatusInternalServerError {
			w.WriteHeader(status)
			return
		}

		if ok := giveMaintenanceResponse(w, region); ok {
			return
		}

		date, expires = profilesCache.AddRecord([]string{region, profileTarget}, data, status)
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
