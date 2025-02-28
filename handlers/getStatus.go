package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bdo-rest-api/cache"
	"bdo-rest-api/config"
	"bdo-rest-api/scraper"
)

var initTime = time.Now()
var version = "1.10.0"

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
				"general":           config.GetCacheTTL().Round(time.Minute).String(),
				"maintenanceStatus": config.GetMaintenanceStatusTTL().Round(time.Minute).String(),
			},
		},
		"docs":      docsLink,
		"proxies":   len(config.GetProxyList()),
		"rateLimit": config.GetRateLimit(),
		"uptime":    time.Since(initTime).Round(time.Second).String(),
		"version":   version,
	})
}
