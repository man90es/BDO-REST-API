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
	"bdo-rest-api/utils"
	"bdo-rest-api/validators"
)

var profileSearchCache = cache.NewCache[[]models.Profile]()

func GetAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	page := validators.ValidatePageQueryParam(r.URL.Query()["page"])
	query, queryOk := validators.ValidateAdventurerNameQueryParam(r.URL.Query()["query"])
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	searchType := validators.ValidateSearchTypeQueryParam(r.URL.Query()["searchType"])

	if !queryOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Look for cached data, then run the scraper if needed
	data, status, date, expires, found := profileSearchCache.GetRecord([]string{region, query, fmt.Sprint(searchType), fmt.Sprint(page)})
	if !found {
		data, status = scrapers.ScrapeAdventurerSearch(region, query, searchType, page)

		if status == http.StatusServiceUnavailable {
			w.Header().Set("Expires", utils.FormatDateForHeaders(time.Now().Add(config.GetMaintenanceStatusTTL())))
			w.WriteHeader(status)
			return
		}

		date, expires = profileSearchCache.AddRecord([]string{region, query, fmt.Sprint(searchType), fmt.Sprint(page)}, data, status)
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
