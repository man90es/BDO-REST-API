package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bdo-rest-api/config"
	"bdo-rest-api/scrapers"
)

var initTime = time.Now()
var version = "1.9.3"

func getStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cache": map[string]interface{}{
			"lastDetectedMaintenance": scrapers.GetLastCloseTimes(),
			"responses": map[string]int{
				"/adventurer":        profilesCache.GetItemCount(),
				"/adventurer/search": profileSearchCache.GetItemCount(),
				"/guild":             guildProfilesCache.GetItemCount(),
				"/guild/search":      guildSearchCache.GetItemCount(),
			},
			"ttl": map[string]string{
				"general":           config.GetCacheTTL().Round(time.Minute).String(),
				"maintenanceStatus": config.GetMaintenanceStatusTTL().Round(time.Minute).String(),
			},
		},
		"docs":    docsLink,
		"proxies": len(config.GetProxyList()),
		"uptime":  time.Since(initTime).Round(time.Second).String(),
		"version": version,
	})
}
