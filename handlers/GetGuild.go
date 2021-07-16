package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scrapers"
)

func GetGuild(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	regionParams, ok1 := r.URL.Query()["region"]
	guildNameParams, ok2 := r.URL.Query()["guildName"]

	if !ok1 || !validateRegion(regionParams[0]) || !ok2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data, status := scrapers.ScrapeGuild(regionParams[0], guildNameParams[0]); status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
