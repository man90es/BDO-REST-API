package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"bdo-rest-api/config"
	"bdo-rest-api/logger"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var scraper *colly.Collector

func init() {
	scraper = colly.NewCollector()
	extensions.RandomUserAgent(scraper)
	scraper.SetRequestTimeout(time.Minute / 2)

	if len(config.GetProxyList()) > 0 {
		scraper.WithTransport(&http.Transport{
			// https://github.com/gocolly/colly/issues/759
			// Make sure that the ProxyFunc is called on every request
			DisableKeepAlives: true,
		})

		scraper.SetProxyFunc(config.GetProxySwitcher())
	}

	scraper.OnError(func(r *colly.Response, err error) {
		logger.Error(fmt.Sprintf("%v", err))
	})

	scraper.OnResponse(func(r *colly.Response) {
		logger.Info(fmt.Sprintf("Received response code for %v: %v", r.Request.URL, r.StatusCode))
	})

	scraper.OnHTML("body", func(body *colly.HTMLElement) {
		if match, _ := regexp.MatchString("/Profile[?]profileTarget=", body.Request.URL.String()); match {
			scrapeAdventurer(body)
		}

		// TODO: Implement scraping for other page types
	})
}

func getRegionPrefix(region string) string {
	return fmt.Sprintf("https://www.%v/Adventure", map[string]string{
		"EU": "naeu.playblackdesert.com/en-US",
		"KR": "kr.playblackdesert.com/ko-KR",
		"SA": "sa.playblackdesert.com/pt-BR",
		"NA": "naeu.playblackdesert.com/en-US",
	}[region])
}

func EnqueueAdventurer(region, profileTarget string) {
	scraper.Visit(fmt.Sprintf("%v/Profile?profileTarget=%v", getRegionPrefix(region), url.QueryEscape(profileTarget)))
}

func EnqueueAdventurerSearch(region, query, searchType string) {
	scraper.Visit(fmt.Sprintf("%v?region=%v&searchType=%v&searchKeyword=%v&Page=1", getRegionPrefix(region), region, searchType, query))
}

func EnqueueGuild(region, name string) {
	scraper.Visit(fmt.Sprintf("%v/Guild/GuildProfile?guildName=%v&region=%v", getRegionPrefix(region), name, region))
}

func EnqueueGuildSearch(region, query string) {
	scraper.Visit(fmt.Sprintf("%v/Guild?region=%v&page=1&searchText=%v", getRegionPrefix(region), region, query))
}
