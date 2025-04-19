package scraper

import (
	"time"

	"github.com/spf13/viper"
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

	expires = lastCloseTimes[region].Add(viper.GetDuration("maintenancettl"))
	return time.Now().Before(expires), expires
}

func setCloseTime(region string) {
	// EU and NA use one website
	if region == "EU" || region == "NA" {
		region = "EUNA"
	}

	lastCloseTimes[region] = time.Now()
}

func GetLastCloseTimes() map[string]time.Time {
	return lastCloseTimes
}
