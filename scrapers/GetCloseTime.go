package scrapers

import (
	"time"

	"bdo-rest-api/config"
)

var lastCloseTimes = map[string]time.Time{
	"EU": {},
	"SA": {},
}

func GetCloseTime(region string) (isCloseTime bool, expires time.Time) {
	// NA and EU use one website
	if region == "NA" {
		region = "EU"
	}

	expires = lastCloseTimes[region].Add(config.GetMaintenanceStatusTTL())
	return time.Now().Before(expires), expires
}

func setCloseTime(region string) {
	lastCloseTimes[region] = time.Now()
}
