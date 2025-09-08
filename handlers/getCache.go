package handlers

import (
	"bdo-rest-api/cache"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

func getCache(w http.ResponseWriter, r *http.Request) {
	if token := viper.GetString("admintoken"); len(token) > 0 {
		providedToken := r.Header.Get("Authorization")

		if providedToken != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	cacheType := r.PathValue("cacheType")

	if len(cacheType) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch cacheType {
	case "adventurer":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cache.Profiles.GetValues())
	case "guild":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cache.GuildProfiles.GetValues())
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
