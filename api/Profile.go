package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"bdo-rest-api/scraper"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	profileTargetParams, ok := r.URL.Query()["profileTarget"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data, status := scraper.ScrapeProfile(url.QueryEscape(profileTargetParams[0])); status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
