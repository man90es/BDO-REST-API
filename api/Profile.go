package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"gitlab.com/man90/black-desert-social-rest-api/scraper"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	profileTargetParams, ok := r.URL.Query()["profileTarget"]

	if !ok {
		return
	}

	if data, err := scraper.ScrapeProfile(url.QueryEscape(profileTargetParams[0])); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(errorResponse{ err.Error() })
	}
}
