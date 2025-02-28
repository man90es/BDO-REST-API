package scraper

import (
	"net/url"
)

func extractProfileTarget(link string) string {
	u, _ := url.Parse(link)
	m, _ := url.ParseQuery(u.RawQuery)
	return m["profileTarget"][0]
}
