package scraper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"

	"bdo-rest-api/logger"
	"bdo-rest-api/utils"
)

func handleTaskError(r *colly.Request, imperva bool, err error) {
	taskRetries, _ := strconv.Atoi(r.Ctx.Get("taskRetries"))

	if imperva {
		logger.Error(fmt.Sprintf("Hit Imperva while loading %v, retries: %v", r.URL, taskRetries))
	} else if strings.Contains(err.Error(), "http2: Transport received GOAWAY from server ErrCode:INTERNAL_ERROR") {
		// This is an error that I don't know how to avoid, it clogs up all future requests
		// and it doesn't seem to be recoverable, so just exit the program and let Docker restart it
		// Maybe in the future there will be a way to avoid/handle this
		logger.Critical(fmt.Sprintf("Error occured: %v", err))
		os.Exit(1)
	} else {
		logger.Error(fmt.Sprintf("Error occured while loading %v: %v, retries: %v", r.URL, err, taskRetries))
	}

	if proxyReloadWebhook := viper.GetString("proxyreloadwebhook"); proxyReloadWebhook != "" {
		utils.SendDummyRequest(proxyReloadWebhook)
	}

	if scraperFailurePause := viper.GetDuration("scraperfailurepause"); scraperFailurePause >= 0 {
		taskQueue.Pause(scraperFailurePause)
	} else {
		taskQueue.Pause(time.Duration(60-time.Now().Second()) * time.Second)
	}

	taskQueue.ConfirmTaskCompletion(r.Ctx.Get("taskClient"), r.Ctx.Get("taskHash"))

	if taskRetries < viper.GetInt("taskretries") {
		taskRegion := r.Ctx.Get("taskRegion")
		taskType := r.Ctx.Get("taskType")
		taskQueue.AddTask(r.Ctx.Get("taskClient"), r.Ctx.Get("taskHash"), utils.BuildRequest(r.URL.String(), map[string]string{
			"taskRegion":  taskRegion,
			"taskRetries": strconv.Itoa(taskRetries + 1),
			"taskType":    taskType,
		}))
	}
}
