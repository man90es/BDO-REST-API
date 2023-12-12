package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/utils"
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

		if isCloseTime, expires := scrapers.GetCloseTime(); isCloseTime {
			w.Header().Set("Expires", utils.FormatDateForHeaders(expires))
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
