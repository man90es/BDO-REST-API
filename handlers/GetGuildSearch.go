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

var guildSearchCache = cache.NewCache[[]models.GuildProfile]()

func GetGuildSearch(w http.ResponseWriter, r *http.Request) {
	name, nameOk := validators.ValidateGuildNameQueryParam(r.URL.Query()["query"])
	page := validators.ValidatePageQueryParam(r.URL.Query()["page"])
	region, regionOk := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !nameOk || !regionOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok := giveMaintenanceResponse(w, region); ok {
		return
	}

	// All names are non-case-sensitive, so this will allow to utilise cache better
	name = strings.ToLower(name)

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := guildSearchCache.GetRecord([]string{region, name, fmt.Sprint(page)})
	if !found {
		data, status = scrapers.ScrapeGuildSearch(region, name, page)

		if status == http.StatusInternalServerError {
			w.WriteHeader(status)
			return
		}

		if ok := giveMaintenanceResponse(w, region); ok {
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
