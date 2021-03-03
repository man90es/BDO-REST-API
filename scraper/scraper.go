package scraper

import (
	"strings"
)

const nice = 69

const closetimeMessage = "BDO servers are currently under maintenance. More info: www.naeu.playblackdesert.com"

func dry(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
