package handlers

import (
	"net/http"
	"strconv"
)

const defaultRegion = "EU"
const defaultPage = 1

func validateRegion(r string) bool {
	return r == "EU" || r == "NA"
}

func validateSearchType(s string) bool {
	return s == "characterName" || s == "familyName"
}

func validatePage(p string) bool {
	page, ok := strconv.Atoi(p)
	return ok == nil && page > 0
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
