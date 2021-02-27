package scraper

import (
	"strings"
)

const nice = 69

func dry(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
