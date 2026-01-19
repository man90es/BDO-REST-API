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

func getAdventurerSearch(w http.ResponseWriter, r *http.Request) {
	region, regionOk, regionValidationMessage := validators.ValidateRegionQueryParam(r.URL.Query()["region"])
	if !regionOk {
		giveBadRequestResponse(w, regionValidationMessage)
		return
	}

	searchTypeQueryParam := r.URL.Query()["searchType"]
	searchType := validators.ValidateSearchTypeQueryParam(searchTypeQueryParam)

	query, queryOk, queryValidationMessage := validators.ValidateAdventurerNameQueryParam(r.URL.Query()["query"], region, searchType)
	if !queryOk {
		giveBadRequestResponse(w, queryValidationMessage)
		return
	}

	bypassCache := validators.ValidateBypassCacheQueryParam(r.URL.Query()["bypassCache"])
	if !bypassCache || !utils.CheckAdminToken(r) {
		if data, status, date, expires, ok := cache.ProfileSearch.GetRecord([]string{region, query, searchType}); !bypassCache && ok {
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

	ok, tasksExceeded, tasksNumber := scraper.EnqueueAdventurerSearch(r.Header.Get("CF-Connecting-IP"), region, query, searchType)
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
			"message": "Player search is being fetched. Please try again later.",
			"status": map[bool]string{
				true:  "started",
				false: "pending",
			}[ok],
		})
	}
}
