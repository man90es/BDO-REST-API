package utils

import (
	"strings"
	"time"
)

func ParseDate(text string) time.Time {
	var format string
	if strings.Contains(text, ".") {
		// Used on KR server for account creation date
		format = "2006.01.02"
	} else if strings.Contains(text, "/") {
		// Used on SA server for account creation date
		format = "02/01/2006"
	} else if strings.Contains(text, "-") {
		// Used on all servers for guild creation date
		format = "2006-01-02"
	} else {
		// Used on NAEU server for account creation date
		format = "Jan 2, 2006"
	}

	if parsed, err := time.Parse(format, RemoveExtraSpaces(text)); nil == err {
		return parsed
	}

	return time.Time{}
}
