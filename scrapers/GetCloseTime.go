package scrapers

import (
	"time"

	"bdo-rest-api/config"
)

var lastCloseTimes = map[string]time.Time{
	"EUNA": {},
	"SA":   {},
}

func GetCloseTime(region string) (isCloseTime bool, expires time.Time) {
	// EU and NA use one website
	if region == "EU" || region == "NA" {
		region = "EUNA"
	}

	expires = lastCloseTimes[region].Add(config.GetMaintenanceStatusTTL())
	return time.Now().Before(expires), expires
}

func setCloseTime(region string) {
	lastCloseTimes[region] = time.Now()
}

func GetLastCloseTimes() map[string]time.Time {
	return lastCloseTimes
}
