package utils

import (
	"net/http"

	"github.com/spf13/viper"
)

func CheckAdminToken(r *http.Request) bool {
	if token := viper.GetString("admintoken"); len(token) > 0 {
		providedToken := r.Header.Get("Authorization")

		if providedToken != "Bearer "+token {
			return false
		}
	}

	return true
}
