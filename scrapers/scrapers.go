package scrapers

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

const closetimeSelector = ".closetime_wrap"

var proxies = make([]string, 0)
var proxySwitcher colly.ProxyFunc

func PushProxies(args ...string) {
	proxies = append(proxies, args...)
	proxySwitcher, _ = proxy.RoundRobinProxySwitcher(proxies...)
}

var verbose = false

func SetVerbose(v bool) {
	verbose = v
}

func dry(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func extractProfileTarget(link string) string {
	u, _ := url.Parse(link)
	m, _ := url.ParseQuery(u.RawQuery)
	return m["profileTarget"][0]
}

func parseDate(text string) time.Time {
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

	if parsed, err := time.Parse(format, dry(text)); nil == err {
		return parsed
	}

	return time.Time{}
}

func getSiteRoot(region string) string {
	if "SA" == region {
		return "https://www.sa.playblackdesert.com/pt-BR"
	}

	if "KR" == region {
		return "https://www.kr.playblackdesert.com/ko-KR"
	}

	return "https://www.naeu.playblackdesert.com/en-US"
}

func collyFactory() (c *colly.Collector) {
	c = colly.NewCollector()
	c.SetRequestTimeout(time.Minute / 2)

	if len(proxies) > 0 {
		c.SetProxyFunc(proxySwitcher)
	}

	c.OnRequest(func(r *colly.Request) {
		if !verbose {
			return
		}

		log.Println("Visiting", r.URL)
	})

	return
}
