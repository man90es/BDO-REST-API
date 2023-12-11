package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bdo-rest-api/cache"
	"bdo-rest-api/config"
	"bdo-rest-api/models"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

var profilesCache = cache.NewCache[models.Profile]()

func GetAdventurer(w http.ResponseWriter, r *http.Request) {
	profileTarget, profileTargetOk := validators.ValidateProfileTargetQueryParam(r.URL.Query()["profileTarget"])
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !profileTargetOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := profilesCache.GetRecord([]string{region, profileTarget})
	if !found {
		data, status = scrapers.ScrapeAdventurer(region, profileTarget)

		if status == http.StatusServiceUnavailable {
			w.Header().Set("Expires", time.Now().Add(config.GetMaintenanceStatusTTL()).Format(time.RFC1123Z))
			w.WriteHeader(status)
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
