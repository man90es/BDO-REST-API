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

	if data, err := scraper.ScrapeProfile(url.QueryEscape(profileTargetParams[0])); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(err.HTTPCode())
		json.NewEncoder(w).Encode(err.Error())
	}
}
