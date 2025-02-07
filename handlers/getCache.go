package handlers

import (
	"bdo-rest-api/cache"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	sf "github.com/sa-/slicefunk"
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
			page, _ := strconv.Atoi(parts[3])

			return map[string]interface{}{
				"region":    parts[0],
				"query":     parts[1],
				"searhType": parts[2],
				"page":      page,
			}
		case "/guild":
			return map[string]interface{}{
				"region":    parts[0],
				"guildName": parts[1],
			}
		case "/guild/search":
			page, _ := strconv.Atoi(parts[2])

			return map[string]interface{}{
				"region": parts[0],
				"query":  parts[1],
				"page":   page,
			}
		default:
			return nil
		}
	}
}

func getCache(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"/adventurer":        sf.Map(cache.Profiles.GetKeys(), getParseCacheKey("/adventurer")),
		"/adventurer/search": sf.Map(cache.ProfileSearch.GetKeys(), getParseCacheKey("/adventurer/search")),
		"/guild":             sf.Map(cache.GuildProfiles.GetKeys(), getParseCacheKey("/guild")),
		"/guild/search":      sf.Map(cache.GuildSearch.GetKeys(), getParseCacheKey("/guild/search")),
	})
}
