package scrapers

import (
	"fmt"
	"log"
	"time"

	"bdo-rest-api/config"

	colly "github.com/gocolly/colly/v2"
)

type scraper struct {
	c      *colly.Collector
	region string
}

func newScraper(region string) (s scraper) {
	s.region = region
	s.c = colly.NewCollector()
	s.c.SetRequestTimeout(time.Minute / 2)

	if len(config.GetProxyList()) > 0 {
		s.c.SetProxyFunc(config.GetProxySwitcher())
	}

	s.c.OnRequest(func(r *colly.Request) {
		if config.GetVerbosity() {
			log.Println("Visiting", r.URL)
		}
	})

	s.OnHTML(`.closetime_wrap`, func(e *colly.HTMLElement) {
		setCloseTime(region)
	})

	return
}

func (s *scraper) OnHTML(goquerySelector string, f colly.HTMLCallback) {
	s.c.OnHTML(goquerySelector, f)
}

func (s *scraper) Visit(URL string) error {
	regionPrefix := map[string]string{
		"EU": "naeu.playblackdesert.com/en-US",
		"KR": "kr.playblackdesert.com/ko-KR",
		"SA": "sa.playblackdesert.com/pt-BR",
		"US": "naeu.playblackdesert.com/en-US",
	}[s.region]

	return s.c.Visit(fmt.Sprintf("https://www.%v/Adventure%v", regionPrefix, URL))
}
