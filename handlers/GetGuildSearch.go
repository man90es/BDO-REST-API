package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/utils"
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

		if isCloseTime, expires := scrapers.GetCloseTime(); isCloseTime {
			w.Header().Set("Expires", utils.FormatDateForHeaders(expires))
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
