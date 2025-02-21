package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/scraper"
	"bdo-rest-api/validators"
)

func getGuild(w http.ResponseWriter, r *http.Request) {
	name, nameOk, nameValidationMessage := validators.ValidateGuildNameQueryParam(r.URL.Query()["guildName"])
	if !nameOk {
		giveBadRequestResponse(w, nameValidationMessage)
		return
	}

	region, regionOk, regionValidationMessage := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	if !regionOk {
		giveBadRequestResponse(w, regionValidationMessage)
		return
	}

	data, status, date, expires, found := cache.GuildProfiles.GetRecord([]string{region, name})
	if !found {
		go scraper.EnqueueGuild(region, name)
		data, status, date, expires = cache.GuildProfiles.WaitForRecord([]string{region, name})
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
