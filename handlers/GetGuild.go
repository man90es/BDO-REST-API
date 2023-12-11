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

var guildProfilesCache = cache.NewCache[models.GuildProfile]()

func GetGuild(w http.ResponseWriter, r *http.Request) {
	name, nameOk := validators.ValidateGuildNameQueryParam(r.URL.Query()["guildName"])
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !nameOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := guildProfilesCache.GetRecord([]string{region, name})
	if !found {
		data, status = scrapers.ScrapeGuild(region, name)

		if status == http.StatusServiceUnavailable {
			w.Header().Set("Expires", time.Now().Add(config.GetMaintenanceStatusTTL()).Format(time.RFC1123Z))
			w.WriteHeader(status)
			return
		}

		date, expires = guildProfilesCache.AddRecord([]string{region, name}, data, status)
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
