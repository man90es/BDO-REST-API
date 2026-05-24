package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// GenSec generates plausible values for Sec-Ch-Ua, Sec-Ch-Ua-Mobile, and Sec-Ch-Ua-Platform headers
func GenSec(useragent string) (ua, mobile, platform string) {
	// 1. Determine Sec-Ch-Ua-Mobile
	mobile = "?0"
	if strings.Contains(useragent, "Mobile") || strings.Contains(useragent, "Android") || strings.Contains(useragent, "iPhone") {
		mobile = "?1"
	}

	// 2. Determine Sec-Ch-Ua-Platform
	platform = "\"Windows\""
	if strings.Contains(useragent, "Android") {
		platform = "\"Android\""
	} else if strings.Contains(useragent, "iPhone") || strings.Contains(useragent, "iPad") {
		platform = "\"iOS\""
	} else if strings.Contains(useragent, "Macintosh") || strings.Contains(useragent, "Mac OS X") {
		platform = "\"macOS\""
	} else if strings.Contains(useragent, "Linux") {
		platform = "\"Linux\""
	}

	// 3. Determine Sec-Ch-Ua (Chromium Brand and Major Version)
	reChrome := regexp.MustCompile(`Chrome/(\d+)`)
	reEdge := regexp.MustCompile(`Edg/(\d+)`)

	if match := reEdge.FindStringSubmatch(useragent); len(match) > 1 {
		ua = fmt.Sprintf("\"Chromium\";v=\"%[1]s\", \"Microsoft Edge\";v=\"%[1]s\", \"Not/A)Brand\";v=\"99\"", match[1])
	} else if match := reChrome.FindStringSubmatch(useragent); len(match) > 1 {
		ua = fmt.Sprintf("\"Chromium\";v=\"%[1]s\", \"Google Chrome\";v=\"%[1]s\", \"Not/A)Brand\";v=\"99\"", match[1])
	} else {
		ua = "\"Chromium\";v=\"131\", \"Google Chrome\";v=\"131\", \"Not/A)Brand\";v=\"99\""
	}

	return
}
