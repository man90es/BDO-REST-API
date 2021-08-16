package scrapers

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

var proxies = make([]string, 0)
var proxySwitcher colly.ProxyFunc

func PushProxies(args ...string) {
	proxies = append(proxies, args...)
	proxySwitcher, _ = proxy.RoundRobinProxySwitcher(proxies...)
}

func dry(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func reEscape(s string) string {
	decodedValue, _ := url.QueryUnescape(s)
	return url.QueryEscape(decodedValue)
}

func extractProfileTarget(link string) string {
	return reEscape(link[69:]) // Nice
}

func collyFactory() (c *colly.Collector) {
	c = colly.NewCollector()
	c.SetRequestTimeout(time.Minute / 2)

	if len(proxies) > 0 {
		c.SetProxyFunc(proxySwitcher)
	}

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	return
}
