package scrapers

import (
	"fmt"
	"time"

	"bdo-rest-api/config"
	"bdo-rest-api/logger"

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

	s.c.OnError(func(r *colly.Response, err error) {
		logger.Error(fmt.Sprintf("%v", err))
	})

	s.c.OnResponse(func(r *colly.Response) {
		logger.Info(fmt.Sprintf("Received response code for %v: %v", r.Request.URL, r.StatusCode))
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
		"NA": "naeu.playblackdesert.com/en-US",
	}[s.region]

	return s.c.Visit(fmt.Sprintf("https://www.%v/Adventure%v", regionPrefix, URL))
}
