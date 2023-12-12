package utils

import "time"

func FormatDateForHeaders(date time.Time) string {
	return date.Format(time.RFC1123Z)
}
