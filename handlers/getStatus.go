package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bdo-rest-api/cache"
	"bdo-rest-api/scraper"

	"github.com/spf13/viper"
)

var initTime = time.Now()
var version = "1.11.3"

func getStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cache": map[string]interface{}{
			"lastDetectedMaintenance": scraper.GetLastCloseTimes(),
			"responses": map[string]int{
				"/adventurer":        cache.Profiles.GetItemCount(),
				"/adventurer/search": cache.ProfileSearch.GetItemCount(),
				"/guild":             cache.GuildProfiles.GetItemCount(),
				"/guild/search":      cache.GuildSearch.GetItemCount(),
			},
			"ttl": map[string]string{
				"general":           viper.GetDuration("cachettl").Round(time.Minute).String(),
				"maintenanceStatus": viper.GetDuration("maintenancettl").Round(time.Minute).String(),
			},
		},
		"docs":      docsLink,
		"proxies":   len(viper.GetStringSlice("proxy")),
		"rateLimit": viper.GetInt64("ratelimit"),
		"taskQueue": map[string]int{
			"maxTasksPerClient": viper.GetInt("maxtasksperclient"),
			"taskRetries":       viper.GetInt("taskretries"),
		},
		"uptime":  time.Since(initTime).Round(time.Second).String(),
		"version": version,
	})
}
