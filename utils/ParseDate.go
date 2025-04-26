package utils

import (
	"strings"
	"time"
)

func ParseDate(text string) time.Time {
	var format string
	if strings.Contains(text, ".") {
		// KR
		format = "2006.01.02"
	} else if strings.Contains(text, "/") {
		// SA
		format = "02/01/2006 (UTC-3)"
	} else {
		// NAEU
		format = "Jan 2, 2006 (UTC)"
	}

	if parsed, err := time.Parse(format, strings.TrimSpace(text)); nil == err {
		return parsed
	}

	return time.Time{}
}
