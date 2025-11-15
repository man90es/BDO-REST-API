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

	"bdo-rest-api/logger"
	"bdo-rest-api/utils"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/spf13/viper"
)

var taskQueue *TaskQueue
var scraperInitialised = false

func InitScraper() {
	if scraperInitialised {
		return
	}
	scraperInitialised = true

	scraper := colly.NewCollector()
	extensions.RandomUserAgent(scraper)
	scraper.AllowURLRevisit = true
	scraper.SetRequestTimeout(time.Minute / 2)

	scraper.Limit(&colly.LimitRule{
		Delay:       time.Second,
		RandomDelay: 5 * time.Second,
	})

	proxies := viper.GetStringSlice("proxy")
	if p, err := proxy.RoundRobinProxySwitcher(proxies...); err == nil {
		scraper.SetProxyFunc(p)
	}

	taskQueue = NewTaskQueue(10000)
	taskQueue.SetProcessFunc(func(t Task) {
		scraper.Visit(t.URL)
	})

	scraper.OnRequest(func(r *colly.Request) {
		query := r.URL.Query()
		for _, key := range []string{"taskHash", "taskType", "taskRegion", "taskRetries", "taskClient"} {
			r.Ctx.Put(key, query.Get(key))
			query.Del(key)
		}
		r.URL.RawQuery = query.Encode()
	})

	scraper.OnError(func(r *colly.Response, err error) {
		logger.Error(fmt.Sprintf("Error occured while loading %v: %v", r.Request.URL, err))
		taskQueue.ConfirmTaskCompletion(r.Ctx.Get("taskClient"), r.Ctx.Get("taskHash"))
	})

	scraper.OnResponse(func(r *colly.Response) {
		logger.Info(fmt.Sprintf("Loaded %v", r.Request.URL))
	})

	scraper.OnHTML("body", func(body *colly.HTMLElement) {
		imperva := false
		queryString, _ := url.ParseQuery(body.Request.URL.RawQuery)
		taskClient := body.Request.Ctx.Get("taskClient")
		taskHash := body.Request.Ctx.Get("taskHash")
		taskRegion := body.Request.Ctx.Get("taskRegion")
		taskType := body.Request.Ctx.Get("taskType")

		body.ForEachWithBreak("iframe", func(_ int, e *colly.HTMLElement) bool {
			imperva = true
			return false
		})

		if imperva {
			taskRetries, _ := strconv.Atoi(body.Request.Ctx.Get("taskRetries"))
			logger.Error(fmt.Sprintf("Hit Imperva while loading %v, retries: %v", body.Request.URL.String(), taskRetries))
			if proxyReloadWebhook := viper.GetString("proxyreloadwebhook"); proxyReloadWebhook != "" {
				utils.SendDummyRequest(proxyReloadWebhook)
				taskQueue.Pause(time.Second * 5)
			} else {
				taskQueue.Pause(time.Duration(60-time.Now().Second()) * time.Second)
			}
			taskQueue.ConfirmTaskCompletion(taskClient, taskHash)

			if taskRetries < viper.GetInt("taskretries") {
				taskQueue.AddTask(taskClient, taskHash, utils.BuildRequest(body.Request.URL.String(), map[string]string{
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
			taskQueue.ConfirmTaskCompletion(taskClient, taskHash)
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

		taskQueue.ConfirmTaskCompletion(taskClient, taskHash)
	})
}

func createTask(clientIP, region, taskType string, query map[string]string) (taskAdded, tasksQuantityExceeded bool) {
	crc32 := crc32.NewIEEE()
	crc32.Write([]byte(strings.Join(append(slices.Sorted(maps.Values(query)), region, taskType), "")))
	hashString := strconv.Itoa(int(crc32.Sum32()))

	if taskQueue.CountQueuedTasksForClient(clientIP) >= viper.GetInt("maxtasksperclient") {
		return false, true
	}

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
		"taskRegion":  region,
		"taskRetries": "0",
		"taskType":    taskType,
	})

	added := taskQueue.AddTask(clientIP, hashString, utils.BuildRequest(url, query))
	return added, false
}

func EnqueueAdventurer(clientIP, region, profileTarget string) (taskAdded, tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "player", map[string]string{
		"profileTarget": profileTarget,
	})
}

func EnqueueAdventurerSearch(clientIP, region, query, searchType string) (taskAdded, tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "playerSearch", map[string]string{
		"Page":          "1",
		"region":        region,
		"searchKeyword": query,
		"searchType":    searchType,
	})
}

func EnqueueGuild(clientIP, region, name string) (taskAdded, tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "guild", map[string]string{
		"guildName": name,
		"region":    region,
	})
}

func EnqueueGuildSearch(clientIP, region, query string) (taskAdded, tasksQuantityExceeded bool) {
	return createTask(clientIP, region, "guildSearch", map[string]string{
		"page":       "1",
		"region":     region,
		"searchText": query,
	})
}
