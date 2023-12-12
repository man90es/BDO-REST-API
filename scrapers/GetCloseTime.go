package scrapers

import (
	"time"

	"bdo-rest-api/config"
)

var lastCloseTime time.Time

func GetCloseTime() (isCloseTime bool, expires time.Time) {
	expires = lastCloseTime.Add(config.GetMaintenanceStatusTTL())
	return time.Now().Before(expires), expires
}

func setCloseTime() {
	lastCloseTime = time.Now()
}
