package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bdo-rest-api/cache"
	"bdo-rest-api/config"
	"bdo-rest-api/models"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

var guildSearchCache = cache.NewCache[[]models.GuildProfile]()

func GetGuildSearch(w http.ResponseWriter, r *http.Request) {
	name, nameOk := validators.ValidateGuildNameQueryParam(r.URL.Query()["query"])
	page := validators.ValidatePageQueryParam(r.URL.Query()["page"])
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !nameOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := guildSearchCache.GetRecord([]string{region, name, fmt.Sprint(page)})
	if !found {
		data, status = scrapers.ScrapeGuildSearch(region, name, page)

		if status == http.StatusServiceUnavailable {
			w.Header().Set("Expires", time.Now().Add(config.GetMaintenanceStatusTTL()).Format(time.RFC1123Z))
			w.WriteHeader(status)
			return
		}

		date, expires = guildSearchCache.AddRecord([]string{region, name, fmt.Sprint(page)}, data, status)
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
