package utils

import "strings"

func RemoveExtraSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
