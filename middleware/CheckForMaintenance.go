package middleware

import (
	"net/http"

	"bdo-rest-api/scrapers"
)

func CheckForMaintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if scrapers.IsCloseTime() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		next.ServeHTTP(w, r)
	})
}
