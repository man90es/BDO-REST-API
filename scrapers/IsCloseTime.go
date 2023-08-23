package scrapers

import (
	"time"
)

var lastCloseTime time.Time

func IsCloseTime() bool {
	return time.Now().Before(lastCloseTime.Add(10 * time.Minute))
}

func setCloseTime() {
	lastCloseTime = time.Now()
}
