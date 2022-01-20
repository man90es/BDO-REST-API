package handlers

import (
	"net/http"
)

const defaultRegion = "EU"
const defaultPage = 1

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
