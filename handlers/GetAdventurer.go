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

func GetAdventurer(w http.ResponseWriter, r *http.Request) {
	profileTarget, profileTargetOk := validators.ValidateProfileTargetQueryParam(r.URL.Query()["profileTarget"])
	region, regionOk := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !profileTargetOk || !regionOk {
		w.WriteHeader(http.StatusBadRequest)
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
