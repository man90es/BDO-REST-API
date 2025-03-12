package handlers

import (
	"net/http"
	"time"

	"bdo-rest-api/scraper"
	"bdo-rest-api/utils"
)

func giveMaintenanceResponse(w http.ResponseWriter, region string) (ok bool) {
	isCloseTime, expires := scraper.GetCloseTime(region)

	if !isCloseTime {
		return false
	}

	w.Header().Set("Date", utils.FormatDateForHeaders(time.Now()))
	w.Header().Set("Expires", utils.FormatDateForHeaders(expires))
	w.WriteHeader(http.StatusServiceUnavailable)
	return true
}
