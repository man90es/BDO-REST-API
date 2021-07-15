package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bdo-rest-api/scraper"
)

func ProfileSearch(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	regionParams, ok1 := r.URL.Query()["region"]
	searchTypeParams, ok2 := r.URL.Query()["searchType"]
	pageParams, ok3 := r.URL.Query()["page"]
	queryParams, ok4 := r.URL.Query()["query"]

	searchType := map[string]int8{
		"characterName": 1,
		"familyName":    2,
	}[searchTypeParams[0]]

	if !ok1 || !ok2 || !ok4 || !validateRegion(regionParams[0]) || searchType < 1 || searchType > 2 {
		return
	}

	page := 1

	if ok3 {
		page, _ = strconv.Atoi(pageParams[0])
	}

	if data, err := scraper.ScrapeProfileSearch(regionParams[0], queryParams[0], searchType, int32(page)); err == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(errorResponse{err.Error()})
	}
}
