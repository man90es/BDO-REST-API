package handlers

import (
	"bdo-rest-api/cache"
	"encoding/json"
	"net/http"
	"strings"

	sf "github.com/sa-/slicefunk"
	"github.com/spf13/viper"
)

func getParseCacheKey(cacheType string) func(string) map[string]interface{} {
	return func(key string) map[string]interface{} {
		parts := strings.Split(key, ",")

		switch cacheType {
		case "/adventurer":
			return map[string]interface{}{
				"region":        parts[0],
				"profileTarget": parts[1],
			}
		case "/adventurer/search":

			return map[string]interface{}{
				"region":    parts[0],
				"query":     parts[1],
				"searhType": parts[2],
			}
		case "/guild":
			return map[string]interface{}{
				"region":    parts[0],
				"guildName": parts[1],
			}
		case "/guild/search":
			return map[string]interface{}{
				"region": parts[0],
				"query":  parts[1],
			}
		default:
			return nil
		}
	}
}

func getCacheSummary(w http.ResponseWriter, r *http.Request) {
	if token := viper.GetString("admintoken"); len(token) > 0 {
		providedToken := r.Header.Get("Authorization")

		if providedToken != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"/adventurer":        sf.Map(cache.Profiles.GetKeys(), getParseCacheKey("/adventurer")),
		"/adventurer/search": sf.Map(cache.ProfileSearch.GetKeys(), getParseCacheKey("/adventurer/search")),
		"/guild":             sf.Map(cache.GuildProfiles.GetKeys(), getParseCacheKey("/guild")),
		"/guild/search":      sf.Map(cache.GuildSearch.GetKeys(), getParseCacheKey("/guild/search")),
	})
}
