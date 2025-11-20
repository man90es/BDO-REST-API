package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/utils"
)

func getCache(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckAdminToken(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
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
