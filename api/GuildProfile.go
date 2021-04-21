package api

import (
	"encoding/json"
	"net/http"

	"gitlab.com/man90/black-desert-social-rest-api/scraper"
)

func GuildProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	regionParams, ok1 := r.URL.Query()["region"]
	guildNameParams, ok2 := r.URL.Query()["guildName"]

	if !ok1 || !validateRegion(regionParams[0]) || !ok2 {
		return
	}

	if data, err := scraper.ScrapeGuildProfile(regionParams[0], guildNameParams[0]); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(errorResponse{err.Error()})
	}
}
