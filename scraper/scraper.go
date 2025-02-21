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
		queryString, _ := url.ParseQuery(body.Request.URL.RawQuery)
		region := queryString["region"][0]

		// TODO: Maintenance detection
		// body.ForEachWithBreak(".type_3", func(_ int, e *colly.HTMLElement) bool {
		// 	setCloseTime(region)
		// 	status = http.StatusServiceUnavailable
		// 	cache.GuildProfiles.SignalMaintenance([]string{region, name}, guildProfile, status)
		// 	return false
		// })

		if match, _ := regexp.MatchString("/Profile[?]profileTarget=", body.Request.URL.String()); match {
			profileTarget := queryString["profileTarget"][0]
			scrapeAdventurer(body, region, profileTarget)
		}

		if match, _ := regexp.MatchString("/Guild/GuildProfile[?]guildName=", body.Request.URL.String()); match {
			scrapeGuild(body, region)
		}

		if match, _ := regexp.MatchString("&searchKeyword=", body.Request.URL.String()); match {
			query := queryString["searchKeyword"][0]
			searchType := queryString["searchType"][0]

			scrapeAdventurerSearch(body, region, query, searchType)
		}

		if match, _ := regexp.MatchString("&page=1&searchText=", body.Request.URL.String()); match {
			query := queryString["searchText"][0]

			scrapeGuildSearch(body, region, query)
		}

		// TODO: Log that none of the scrapers fits
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
	scraper.Visit(fmt.Sprintf("%v/Profile?profileTarget=%v&region=%v", getRegionPrefix(region), url.QueryEscape(profileTarget), region))
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
