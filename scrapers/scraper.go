package scrapers

import (
	"fmt"
	"log"
	"time"

	colly "github.com/gocolly/colly/v2"
)

type scraper struct {
	c *colly.Collector
}

func newScraper() (s scraper) {
	s.c = colly.NewCollector()
	s.c.SetRequestTimeout(time.Minute / 2)

	if len(proxies) > 0 {
		s.c.SetProxyFunc(proxySwitcher)
	}

	s.c.OnRequest(func(r *colly.Request) {
		if isVerbose() {
			log.Println("Visiting", r.URL)
		}
	})

	s.OnHTML(`.closetime_wrap`, func(e *colly.HTMLElement) {
		setCloseTime()
	})

	return
}

func (s *scraper) OnHTML(goquerySelector string, f colly.HTMLCallback) {
	s.c.OnHTML(goquerySelector, f)
}

func (s *scraper) Visit(URL string, region string) error {
	regionPrefix := map[string]string{
		"EU": "naeu.playblackdesert.com/en-US",
		"KR": "kr.playblackdesert.com/ko-KR",
		"SA": "sa.playblackdesert.com/pt-BR",
		"US": "naeu.playblackdesert.com/en-US",
	}[region]

	return s.c.Visit(fmt.Sprintf("https://www.%v/Adventure%v", regionPrefix, URL))
}
