package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"bdo-rest-api/scrapers"
)

func GetAdventurer(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	profileTargetParams, ok := r.URL.Query()["profileTarget"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data, status := scrapers.ScrapeAdventurer(url.QueryEscape(profileTargetParams[0])); status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
