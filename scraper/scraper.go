package scraper

import (
	"fmt"
	"hash/crc32"
	"maps"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"bdo-rest-api/config"
	"bdo-rest-api/logger"
	"bdo-rest-api/utils"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
)

var taskQueue *TaskQueue

func init() {
	scraper := colly.NewCollector()
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

	taskQueue = NewTaskQueue(10000)
	taskQueue.SetProcessFunc(func(t Task) {
		scraper.Visit(t.URL)
	})

	scraper.OnRequest(func(r *colly.Request) {
		query := r.URL.Query()

		r.Ctx.Put("taskHash", query.Get("taskHash"))
		query.Del("taskHash")

		r.Ctx.Put("taskType", query.Get("taskType"))
		query.Del("taskType")

		r.Ctx.Put("taskRegion", query.Get("taskRegion"))
		query.Del("taskRegion")

		r.Ctx.Put("taskRetries", query.Get("taskRetries"))
		query.Del("taskRetries")

		r.Ctx.Put("taskClient", query.Get("taskClient"))
		query.Del("taskClient")

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
		taskHash := body.Request.Ctx.Get("taskHash")
		taskRegion := body.Request.Ctx.Get("taskRegion")
		taskType := body.Request.Ctx.Get("taskType")

		body.ForEachWithBreak("iframe", func(_ int, e *colly.HTMLElement) bool {
			imperva = true
			return false
		})

		if imperva {
			taskRetries, _ := strconv.Atoi(body.Request.Ctx.Get("taskRetries"))
			taskClient := body.Request.Ctx.Get("taskClient")

			logger.Error(fmt.Sprintf("Hit Imperva while loading %v, retries: %v", body.Request.URL.String(), taskRetries))
			taskQueue.Pause(time.Duration(60-time.Now().Second()) * time.Second)

			// TODO: Make this configurable
			if taskRetries < 3 {
				taskHashI, _ := strconv.Atoi(taskHash)
				taskQueue.AddTask(taskClient, uint32(taskHashI), utils.BuildRequest(body.Request.URL.String(), map[string]string{
					"taskClient":  taskClient,
					"taskHash":    taskHash,
					"taskRegion":  taskRegion,
					"taskRetries": strconv.Itoa(taskRetries + 1),
					"taskType":    taskType,
				}))
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

func createTask(clientIP, region, taskType string, query map[string]string) (tasksQuantityExceeded bool) {
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

	crc32 := crc32.NewIEEE()
	crc32.Write([]byte(strings.Join(append(slices.Collect(maps.Values(query)), url), "")))

	if unique := taskQueue.CheckHashUnique(crc32.Sum32()); !unique {
		return false
	}

	// TODO: Make this configurable
	if taskQueue.CountQueuedTasksForClient(clientIP) >= 5 {
		return true
	}

	maps.Copy(query, map[string]string{
		"taskClient":  clientIP,
		"taskHash":    strconv.Itoa(int(crc32.Sum32())),
		"taskRegion":  region,
		"taskRetries": "0",
		"taskType":    taskType,
	})

	taskQueue.AddTask(clientIP, crc32.Sum32(), utils.BuildRequest(url, query))
	return false
}

func EnqueueAdventurer(clientIP, region, profileTarget string) (tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "player", map[string]string{
		"profileTarget": profileTarget,
	})
}

func EnqueueAdventurerSearch(clientIP, region, query, searchType string) (tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "playerSearch", map[string]string{
		"Page":          "1",
		"region":        region,
		"searchKeyword": query,
		"searchType":    searchType,
	})
}

func EnqueueGuild(clientIP, region, name string) (tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "guild", map[string]string{
		"guildName": name,
		"region":    region,
	})
}

func EnqueueGuildSearch(clientIP, region, query string) (tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "guildSearch", map[string]string{
		"page":       "1",
		"region":     region,
		"searchText": query,
	})
}
