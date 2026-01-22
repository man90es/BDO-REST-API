package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bdo-rest-api/cache"
	"bdo-rest-api/scraper"
	"bdo-rest-api/utils"
	"bdo-rest-api/validators"
)

func getAdventurer(w http.ResponseWriter, r *http.Request) {
	profileTarget, profileTargetOk, profileTargetValidationMessage := validators.ValidateProfileTargetQueryParam(r.URL.Query()["profileTarget"])
	if !profileTargetOk {
		giveBadRequestResponse(w, profileTargetValidationMessage)
		return
	}

	region, regionOk, regionValidationMessage := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	if !regionOk {
		giveBadRequestResponse(w, regionValidationMessage)
		return
	}

	bypassCache := validators.ValidateBypassCacheQueryParam(r.URL.Query()["bypassCache"])
	if !bypassCache || !utils.CheckAdminToken(r) {
		if data, status, date, expires, ok := cache.Profiles.GetRecord([]string{region, profileTarget}); ok {
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

	ok, tasksExceeded, tasksNumber := scraper.EnqueueAdventurer(r.Header.Get("CF-Connecting-IP"), region, profileTarget)
	w.Header().Set("X-Tasks-Number", strconv.Itoa(tasksNumber))
	if tasksExceeded {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "You have exceeded the maximum number of concurrent tasks.",
			"status":  "rejected",
		})
	} else {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Player profile is being fetched. Please try again later.",
			"status": map[bool]string{
				true:  "started",
				false: "pending",
			}[ok],
		})
	}
}
