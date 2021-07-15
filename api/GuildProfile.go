package api

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/scraper"
)

func GuildProfile(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	regionParams, ok1 := r.URL.Query()["region"]
	guildNameParams, ok2 := r.URL.Query()["guildName"]

	if !ok1 || !validateRegion(regionParams[0]) || !ok2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data, err := scraper.ScrapeGuildProfile(regionParams[0], guildNameParams[0]); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(err.HTTPCode())
		json.NewEncoder(w).Encode(err.Error())
	}
}
