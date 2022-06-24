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

func dry(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func extractProfileTarget(link string) string {
	u, _ := url.Parse(link)
	m, _ := url.ParseQuery(u.RawQuery)
	return url.QueryEscape(m["profileTarget"][0])
}

func getSiteRoot(region string) string {
	if "SA" == region {
		return "https://www.sa.playblackdesert.com/pt-BR"
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
		log.Println("Visiting", r.URL)
	})

	return
}
