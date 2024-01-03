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
	useragent := colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	s.region = region
	s.c = colly.NewCollector(useragent)
	s.c.SetRequestTimeout(time.Minute / 2)

	if len(config.GetProxyList()) > 0 {
		s.c.SetProxyFunc(config.GetProxySwitcher())
	}

	s.c.OnRequest(func(r *colly.Request) {
		if config.GetVerbosity() {
			log.Println("Visiting", r.URL)
		}
	})

	// Detect maintenance
	s.OnHTML(`.type_3`, func(e *colly.HTMLElement) {
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
