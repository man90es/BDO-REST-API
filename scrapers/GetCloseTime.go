package scrapers

import (
	"time"

	"bdo-rest-api/config"
)

var lastCloseTimes = map[string]time.Time{
	"EU": time.Time{},
	"NA": time.Time{},
	"SA": time.Time{},
}

func GetCloseTime(region string) (isCloseTime bool, expires time.Time) {
	expires = lastCloseTimes[region].Add(config.GetMaintenanceStatusTTL())
	return time.Now().Before(expires), expires
}

func setCloseTime(region string) {
	lastCloseTimes[region] = time.Now()
}
