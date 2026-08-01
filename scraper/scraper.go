package scraper

import (
	"fmt"
	"hash/crc32"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"bdo-rest-api/logger"
	"bdo-rest-api/utils"

	colly "github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	ua "github.com/nzrsky/useragent-generator/pkg/useragent"
	"github.com/spf13/viper"
)

var taskQueue *TaskQueue
var scraperInitialised = false

const metadataTaskAddedAt = "taskAddedAt"
const metadataTaskClient = "taskClient"
const metadataTaskHash = "taskHash"
const metadataTaskRegion = "taskRegion"
const metadataTaskRetries = "taskRetries"
const metadataTaskType = "taskType"

func InitScraper() {
	if scraperInitialised {
		return
	}
	scraperInitialised = true

	scraper := colly.NewCollector()
	scraper.AllowURLRevisit = true
	scraper.SetRequestTimeout(time.Minute / 2)

	proxies := viper.GetStringSlice("proxy")
	if p, err := proxy.RoundRobinProxySwitcher(proxies...); err == nil {
		scraper.SetProxyFunc(p)
	}

	taskQueue = NewTaskQueue(10000)
	taskQueue.SetProcessFunc(func(t Task) {
		ctx := colly.NewContext()
		for k, v := range t.Metadata {
			ctx.Put(k, v)
		}

		u, _ := url.Parse(t.URL)
		useragent := ua.Random()
		secUa, secMobile, secPlatform := utils.GenSec(useragent)
		hdr := make(http.Header)
		hdr.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		hdr.Set("Accept-Encoding", "gzip, deflate, br, zstd")
		hdr.Set("Accept-Language", "en-US,en;q=0.9")
		hdr.Set("Cache-Control", "max-age=0")
		hdr.Set("Priority", "u=0, i")
		hdr.Set("Referer", fmt.Sprintf("https://%v/%v/Main/Index", u.Host, u.Path[1:6]))
		hdr.Set("Sec-Ch-Ua", secUa)
		hdr.Set("Sec-Ch-Ua-Mobile", secMobile)
		hdr.Set("Sec-Ch-Ua-Platform", secPlatform)
		hdr.Set("Sec-Fetch-Dest", "document")
		hdr.Set("Sec-Fetch-Mode", "navigate")
		hdr.Set("Sec-Fetch-Site", "same-origin")
		hdr.Set("Sec-Fetch-User", "?1")
		hdr.Set("Upgrade-Insecure-Requests", "1")
		hdr.Set("User-Agent", useragent)

		scraper.Request("GET", t.URL, nil, ctx, hdr)
	})

	scraper.OnError(func(r *colly.Response, err error) {
		handleTaskError(r.Request, false, err)
	})

	scraper.OnHTML("body", func(body *colly.HTMLElement) {
		queryString, _ := url.ParseQuery(body.Request.URL.RawQuery)
		taskClient := body.Request.Ctx.Get(metadataTaskClient)
		taskHash := body.Request.Ctx.Get(metadataTaskHash)
		taskRegion := body.Request.Ctx.Get(metadataTaskRegion)
		taskType := body.Request.Ctx.Get(metadataTaskType)

		if strings.Contains(body.Text, "Incapsula incident ID") {
			handleTaskError(body.Request, true, nil)
			return
		}

		parsedStartTime, _ := time.Parse(time.RFC3339, body.Request.Ctx.Get(metadataTaskAddedAt))
		elapsed := time.Since(parsedStartTime)
		logger.Info(fmt.Sprintf("Loaded %v in %v", body.Request.URL, elapsed))

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

func createTask(taskClient, region, taskType string, query map[string]string) (ok, tasksExceeded bool, tasksNumber int) {
	crc32 := crc32.NewIEEE()
	crc32.Write([]byte(strings.Join(append(slices.Sorted(maps.Values(query)), region, taskType), "")))
	hashString := strconv.Itoa(int(crc32.Sum32()))

	tasksN := taskQueue.CountQueuedTasksForClient(taskClient)
	if tasksN >= viper.GetInt("maxtasksperclient") {
		return false, true, tasksN
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

	ok = taskQueue.AddTask(
		taskClient,
		hashString,
		utils.BuildRequest(url, query),
		false,
		map[string]string{
			metadataTaskAddedAt: time.Now().Format(time.RFC3339),
			metadataTaskClient:  taskClient,
			metadataTaskHash:    hashString,
			metadataTaskRegion:  region,
			metadataTaskRetries: "0",
			metadataTaskType:    taskType,
		},
	)

	return ok, false, map[bool]int{true: tasksN + 1, false: tasksN}[ok]
}

func EnqueueAdventurer(taskClient, region, profileTarget string) (ok, tasksExceeded bool, tasksNumber int) {
	return createTask(taskClient, region, "player", map[string]string{
		"profileTarget": profileTarget,
	})
}

func EnqueueAdventurerSearch(taskClient, region, query, searchType string) (ok, tasksExceeded bool, tasksNumber int) {
	return createTask(taskClient, region, "playerSearch", map[string]string{
		"Page":          "1",
		"region":        region,
		"searchKeyword": query,
		"searchType":    searchType,
	})
}

func EnqueueGuild(taskClient, region, name string) (ok, tasksExceeded bool, tasksNumber int) {
	return createTask(taskClient, region, "guild", map[string]string{
		"guildName": name,
		"region":    region,
	})
}

func EnqueueGuildSearch(taskClient, region, query string) (ok, tasksExceeded bool, tasksNumber int) {
	return createTask(taskClient, region, "guildSearch", map[string]string{
		"page":       "1",
		"region":     region,
		"searchText": query,
	})
}
