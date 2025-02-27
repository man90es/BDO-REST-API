package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bdo-rest-api/cache"
	"bdo-rest-api/config"
	"bdo-rest-api/logger"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/google/uuid"
)

var scraper *colly.Collector

func init() {
	scraper = colly.NewCollector()
	extensions.RandomUserAgent(scraper)
	scraper.AllowURLRevisit = true
	scraper.SetRequestTimeout(time.Minute / 2)

	scraper.Limit(&colly.LimitRule{
		Delay:       time.Second,
		RandomDelay: 5 * time.Second,
	})

	if len(config.GetProxyList()) > 0 {
		scraper.WithTransport(&http.Transport{
			// https://github.com/gocolly/colly/issues/759
			// Make sure that the ProxyFunc is called on every request
			DisableKeepAlives: true,
		})

		scraper.SetProxyFunc(config.GetProxySwitcher())
	}

	scraper.OnRequest(func(r *colly.Request) {
		query := r.URL.Query()

		r.Ctx.Put("taskId", query.Get("taskId"))
		query.Del("taskId")

		r.Ctx.Put("taskType", query.Get("taskType"))
		query.Del("taskType")

		r.Ctx.Put("taskRegion", query.Get("taskRegion"))
		query.Del("taskRegion")

		r.Ctx.Put("taskRetries", query.Get("taskRetries"))
		query.Del("taskRetries")

		r.URL.RawQuery = query.Encode()
	})

	scraper.OnError(func(r *colly.Response, err error) {
		logger.Error(fmt.Sprintf("Error occured while loading %v: %v", r.Request.URL, err))
	})

	scraper.OnResponse(func(r *colly.Response) {
		logger.Info(fmt.Sprintf("Loaded %v", r.Request.URL))
	})

	scraper.OnHTML("body", func(body *colly.HTMLElement) {
		imperva := false
		queryString, _ := url.ParseQuery(body.Request.URL.RawQuery)
		region := body.Request.Ctx.Get("taskRegion")

		body.ForEachWithBreak("iframe", func(_ int, e *colly.HTMLElement) bool {
			imperva = true
			return false
		})

		if imperva {
			logger.Error("Imperva")
			time.Sleep(10 * time.Second)

			retries, _ := strconv.Atoi(body.Request.Ctx.Get("taskRetries"))
			taskType := body.Request.Ctx.Get("taskType")

			if retries < 2 {
				url := fmt.Sprintf(
					"%v&taskId=%v&taskType=%v&taskRegion=%v&taskRetries=%v",
					body.Request.URL.String(),
					body.Request.Ctx.Get("taskId"),
					taskType,
					body.Request.Ctx.Get("taskRegion"),
					retries+1,
				)
				go scraper.Visit(url)
			} else {
				switch taskType {
				case "player":
					cache.Profiles.SignalBypassCache(http.StatusInternalServerError, body.Request.Ctx.Get("taskId"))

				case "playerSearch":
					cache.ProfileSearch.SignalBypassCache(http.StatusInternalServerError, body.Request.Ctx.Get("taskId"))

				case "guild":
					cache.GuildProfiles.SignalBypassCache(http.StatusInternalServerError, body.Request.Ctx.Get("taskId"))

				case "guildSearch":
					cache.GuildSearch.SignalBypassCache(http.StatusInternalServerError, body.Request.Ctx.Get("taskId"))

				default:
					logger.Error(fmt.Sprintf("Task type %v doesn't match any defined error handlers", taskType))
				}
			}

			return
		}

		body.ForEachWithBreak(".type_3", func(_ int, e *colly.HTMLElement) bool {
			// Request gets redirected to https://www.naeu.playblackdesert.com/en-US/shutdown/closetime?shutDownType=0
			// Maybe a better way to detect maintenance would be looking at the URL
			setCloseTime(region)
			return false
		})

		if isCloseTime, _ := GetCloseTime(region); isCloseTime {
			switch body.Request.Ctx.Get("taskType") {
			case "player":
				cache.Profiles.SignalBypassCache(http.StatusServiceUnavailable, body.Request.Ctx.Get("taskId"))

			case "playerSearch":
				cache.ProfileSearch.SignalBypassCache(http.StatusServiceUnavailable, body.Request.Ctx.Get("taskId"))

			case "guild":
				cache.GuildProfiles.SignalBypassCache(http.StatusServiceUnavailable, body.Request.Ctx.Get("taskId"))

			case "guildSearch":
				cache.GuildSearch.SignalBypassCache(http.StatusServiceUnavailable, body.Request.Ctx.Get("taskId"))

			default:
				logger.Error(fmt.Sprintf("Task type %v doesn't match any defined maintenance handlers", body.Request.Ctx.Get("taskType")))
			}

			return
		}

		switch body.Request.Ctx.Get("taskType") {
		case "player":
			profileTarget := queryString["profileTarget"][0]
			scrapeAdventurer(body, region, profileTarget)

		case "playerSearch":
			query := queryString["searchKeyword"][0]
			searchType := queryString["searchType"][0]
			scrapeAdventurerSearch(body, region, query, searchType)

		case "guild":
			guildName := queryString["guildName"][0]
			scrapeGuild(body, region, guildName)

		case "guildSearch":
			query := queryString["searchText"][0]
			scrapeGuildSearch(body, region, query)

		default:
			logger.Error(fmt.Sprintf("Task type %v doesn't match any defined scrapers", body.Request.Ctx.Get("taskType")))
		}

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

func EnqueueAdventurer(region, profileTarget string) (taskId string, maintenance bool) {
	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		maintenance = true
		return
	}

	taskId = uuid.New().String()
	url := fmt.Sprintf(
		"%v/Profile?profileTarget=%v&taskId=%v&taskType=player&taskRegion=%v&taskRetries=%v",
		getRegionPrefix(region),
		url.QueryEscape(profileTarget),
		taskId,
		region,
		0,
	)
	go scraper.Visit(url)

	return
}

func EnqueueAdventurerSearch(region, query, searchType string) (taskId string, maintenance bool) {
	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		maintenance = true
		return
	}

	taskId = uuid.New().String()
	url := fmt.Sprintf(
		"%v?region=%v&searchType=%v&searchKeyword=%v&Page=1&taskId=%v&taskType=playerSearch&taskRegion=%v&taskRetries=%v",
		getRegionPrefix(region),
		region,
		searchType,
		query,
		taskId,
		region,
		0,
	)
	go scraper.Visit(url)

	return
}

func EnqueueGuild(region, name string) (taskId string, maintenance bool) {
	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		maintenance = true
		return
	}

	taskId = uuid.New().String()
	url := fmt.Sprintf(
		"%v/Guild/GuildProfile?guildName=%v&region=%v&taskId=%v&taskType=guild&taskRegion=%v&taskRetries=%v",
		getRegionPrefix(region),
		name,
		region,
		taskId,
		region,
		0,
	)
	go scraper.Visit(url)

	return
}

func EnqueueGuildSearch(region, query string) (taskId string, maintenance bool) {
	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		maintenance = true
		return
	}

	taskId = uuid.New().String()
	url := fmt.Sprintf(
		"%v/Guild?region=%v&page=1&searchText=%v&taskId=%v&taskType=guildSearch&taskRegion=%v&taskRetries=%v",
		getRegionPrefix(region),
		region,
		query,
		taskId,
		region,
		0,
	)
	go scraper.Visit(url)

	return
}
