package utils

import "net/url"

func BuildRequest(urlString string, queryMap map[string]string) string {
	url, _ := url.Parse(urlString)
	q := url.Query()

	for key, value := range queryMap {
		q.Set(key, value)
	}

	url.RawQuery = q.Encode()
	return url.String()
}
