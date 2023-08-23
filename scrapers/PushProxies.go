package scrapers

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

var proxies = make([]string, 0)
var proxySwitcher colly.ProxyFunc

func PushProxies(args ...string) {
	proxies = append(proxies, args...)
	proxySwitcher, _ = proxy.RoundRobinProxySwitcher(proxies...)
}
