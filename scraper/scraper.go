package scraper

import (
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bdo-rest-api/cache"
	"bdo-rest-api/config"
	"bdo-rest-api/logger"
	"bdo-rest-api/utils"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
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

	if p, err := proxy.RoundRobinProxySwitcher(config.GetProxyList()...); err == nil {
		scraper.SetProxyFunc(p)
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
		taskRegion := body.Request.Ctx.Get("taskRegion")
		taskRetries, _ := strconv.Atoi(body.Request.Ctx.Get("taskRetries"))
		taskId := body.Request.Ctx.Get("taskId")
		taskType := body.Request.Ctx.Get("taskType")

		body.ForEachWithBreak("iframe", func(_ int, e *colly.HTMLElement) bool {
			imperva = true
			return false
		})

		if imperva {
			logger.Error("Imperva")
			time.Sleep(10 * time.Second)

			if taskRetries < 3 {
				go scraper.Visit(utils.BuildRequest(body.Request.URL.String(), map[string]string{
					"taskId":      taskId,
					"taskRegion":  taskRegion,
					"taskRetries": strconv.Itoa(taskRetries + 1),
					"taskType":    taskType,
				}))
			} else {
				switch taskType {
				case "player":
					cache.Profiles.SignalBypassCache(http.StatusInternalServerError, taskId)

				case "playerSearch":
					cache.ProfileSearch.SignalBypassCache(http.StatusInternalServerError, taskId)

				case "guild":
					cache.GuildProfiles.SignalBypassCache(http.StatusInternalServerError, taskId)

				case "guildSearch":
					cache.GuildSearch.SignalBypassCache(http.StatusInternalServerError, taskId)

				default:
					logger.Error(fmt.Sprintf("Task type %v doesn't match any defined error handlers", taskType))
				}
			}

			return
		}

		body.ForEachWithBreak(".type_3", func(_ int, e *colly.HTMLElement) bool {
			// Request gets redirected to https://www.naeu.playblackdesert.com/en-US/shutdown/closetime?shutDownType=0
			// Maybe a better way to detect maintenance would be looking at the URL
			setCloseTime(taskRegion)
			return false
		})

		if isCloseTime, _ := GetCloseTime(taskRegion); isCloseTime {
			switch taskType {
			case "player":
				cache.Profiles.SignalBypassCache(http.StatusServiceUnavailable, taskId)

			case "playerSearch":
				cache.ProfileSearch.SignalBypassCache(http.StatusServiceUnavailable, taskId)

			case "guild":
				cache.GuildProfiles.SignalBypassCache(http.StatusServiceUnavailable, taskId)

			case "guildSearch":
				cache.GuildSearch.SignalBypassCache(http.StatusServiceUnavailable, taskId)

			default:
				logger.Error(fmt.Sprintf("Task type %v doesn't match any defined maintenance handlers", taskType))
			}

			return
		}

		switch taskType {
		case "player":
			profileTarget := queryString["profileTarget"][0]
			scrapeAdventurer(body, taskRegion, profileTarget)

		case "playerSearch":
			query := queryString["searchKeyword"][0]
			searchType := queryString["searchType"][0]
			scrapeAdventurerSearch(body, taskRegion, query, searchType)

		case "guild":
			guildName := queryString["guildName"][0]
			scrapeGuild(body, taskRegion, guildName)

		case "guildSearch":
			query := queryString["searchText"][0]
			scrapeGuildSearch(body, taskRegion, query)

		default:
			logger.Error(fmt.Sprintf("Task type %v doesn't match any defined scrapers", taskType))
		}
	})
}

func createTask(region, taskType string, query map[string]string) (taskId string, maintenance bool) {
	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		return "", true
	}

	taskId = uuid.New().String()

	url := fmt.Sprintf(
		"https://www.%v/Adventure%v",
		map[string]string{
			"EU": "naeu.playblackdesert.com/en-US",
			"KR": "kr.playblackdesert.com/ko-KR",
			"SA": "sa.playblackdesert.com/pt-BR",
			"NA": "naeu.playblackdesert.com/en-US",
		}[region],
		map[string]string{
			"guild":        "/Guild/GuildProfile",
			"guildSearch":  "/Guild",
			"player":       "/Profile",
			"playerSearch": "",
		}[taskType],
	)

	maps.Copy(query, map[string]string{
		"taskId":      taskId,
		"taskRegion":  region,
		"taskRetries": "0",
		"taskType":    taskType,
	})

	go scraper.Visit(utils.BuildRequest(url, query))
	return
}

func EnqueueAdventurer(region, profileTarget string) (taskId string, maintenance bool) {
	return createTask(region, "player", map[string]string{
		"profileTarget": profileTarget,
	})
}

func EnqueueAdventurerSearch(region, query, searchType string) (taskId string, maintenance bool) {
	return createTask(region, "playerSearch", map[string]string{
		"Page":          "1",
		"region":        region,
		"searchKeyword": query,
		"searchType":    searchType,
	})
}

func EnqueueGuild(region, name string) (taskId string, maintenance bool) {
	return createTask(region, "guild", map[string]string{
		"guildName": name,
		"region":    region,
	})
}

func EnqueueGuildSearch(region, query string) (taskId string, maintenance bool) {
	return createTask(region, "guildSearch", map[string]string{
		"page":       "1",
		"region":     region,
		"searchText": query,
	})
}
