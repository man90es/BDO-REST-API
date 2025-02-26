package handlers

import (
	"encoding/json"
	"net/http"

	"bdo-rest-api/cache"
	"bdo-rest-api/scraper"
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

	data, status, date, expires, found := cache.Profiles.GetRecord([]string{region, profileTarget})
	if !found {
		taskId := scraper.EnqueueAdventurer(region, profileTarget)
		data, status, date, expires = cache.Profiles.WaitForRecord(taskId)
	}

	w.Header().Set("Date", date)
	w.Header().Set("Expires", expires)

	if status == http.StatusOK {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}
