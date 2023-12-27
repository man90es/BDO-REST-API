package middleware

import (
	"net/http"
	"time"
)

func SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Date", time.Now().Format(time.RFC1123Z))

		next.ServeHTTP(w, r)
	})
}
