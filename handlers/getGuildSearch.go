package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func getGuildSearch(w http.ResponseWriter, r *http.Request) {
	name, nameOk, nameValidationMessage := validators.ValidateGuildNameQueryParam(r.URL.Query()["query"])
	if !nameOk {
		giveBadRequestResponse(w, nameValidationMessage)
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
	data, status, date, expires, found := cache.GuildSearch.GetRecord([]string{region, name})
	if !found {
		data, status, date, expires = scrapers.ScrapeGuildSearch(region, name)

		if ok := giveMaintenanceResponse(w, region); ok {
			return
		}
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
