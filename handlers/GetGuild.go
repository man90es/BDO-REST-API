package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/validators"
)

func GetGuild(w http.ResponseWriter, r *http.Request) {
	name, nameOk := validators.ValidateGuildNameQueryParam(r.URL.Query()["guildName"])
	region := validators.ValidateRegionQueryParam(r.URL.Query()["region"])

	if !nameOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, status := scrapers.ScrapeGuild(region, name)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
