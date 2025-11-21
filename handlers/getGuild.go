package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/scraper"
	"bdo-rest-api/utils"
	"bdo-rest-api/validators"
)

func getGuild(w http.ResponseWriter, r *http.Request) {
	name, nameOk, nameValidationMessage := validators.ValidateGuildNameQueryParam(r.URL.Query()["guildName"])
	if !nameOk {
		giveBadRequestResponse(w, nameValidationMessage)
		return
	}

	region, regionOk, regionValidationMessage := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	if !regionOk {
		giveBadRequestResponse(w, regionValidationMessage)
		return
	}

	bypassCache := validators.ValidateBypassCacheQueryParam(r.URL.Query()["bypassCache"])
	if !bypassCache || !utils.CheckAdminToken(r) {
		if data, status, date, expires, ok := cache.GuildProfiles.GetRecord([]string{region, name}); !bypassCache && ok {
			w.Header().Set("Expires", expires)
			w.Header().Set("Last-Modified", date)

			if status == http.StatusOK {
				json.NewEncoder(w).Encode(data)
			} else {
				w.WriteHeader(status)
			}

			return
		}
	}

	if ok := giveMaintenanceResponse(w, region); ok {
		return
	}

	if taskAdded, tasksQuantityExceeded := scraper.EnqueueGuild(r.Header.Get("CF-Connecting-IP"), region, name); tasksQuantityExceeded {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "You have exceeded the maximum number of concurrent tasks.",
			"status":  "rejected",
		})
	} else {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Guild profile is being fetched. Please try again later.",
			"status": map[bool]string{
				true:  "started",
				false: "pending",
			}[taskAdded],
		})
	}
}
