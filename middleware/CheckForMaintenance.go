package middleware

import (
	"net/http"

	"bdo-rest-api/scrapers"
	"bdo-rest-api/utils"
)

func CheckForMaintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isCloseTime, expires := scrapers.GetCloseTime(); isCloseTime {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Expires", utils.FormatDateForHeaders(expires))
			return
		}

		next.ServeHTTP(w, r)
	})
}
