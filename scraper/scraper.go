package scraper

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

const closetimeMessage = "BDO servers are currently under maintenance. More info: www.naeu.playblackdesert.com"

var proxies = make([]string, 0)
var proxySwitcher colly.ProxyFunc

func PushProxies(args ...string) {
	proxies = append(proxies, args...)
	proxySwitcher, _ = proxy.RoundRobinProxySwitcher(proxies...)
}

func dry(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func extractProfileTarget(link string) string {
	return link[69:] // Nice
}

func collyFactory() (c *colly.Collector) {
	c = colly.NewCollector()
	c.SetRequestTimeout(time.Minute / 2)

	if len(proxies) > 0 {
		c.SetProxyFunc(proxySwitcher)
	}

	return
}
