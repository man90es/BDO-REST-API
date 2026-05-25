package scraper

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/utils"
)

func scrapeGuildSearch(body *colly.HTMLElement, region, query string) {
	status := http.StatusNotFound
	guildProfiles := make([]models.GuildProfile, 0)

	body.ForEach(".box_list_area li:not(.no_result)", func(_ int, e *colly.HTMLElement) {
		createdOn := utils.ParseDate(e.ChildText(".date"))
		status = http.StatusOK

		guildProfile := models.GuildProfile{
			Name:   e.ChildText(".guild_title a"),
			Region: region,
			Master: &models.Profile{
				FamilyName:    e.ChildText(".guild_info a"),
				ProfileTarget: extractProfileTarget(e.ChildAttr(".guild_info a", "href")),
			},
			CreatedOn: &createdOn,
		}

		if membersStr := e.ChildText(".member"); true {
			population, _ := strconv.Atoi(membersStr)
			guildProfile.Population = uint8(population)
		}

		guildProfiles = append(guildProfiles, guildProfile)
	})

	startTime, _ := time.Parse(time.RFC3339, body.Request.Ctx.Get(metadataTaskAddedAt))
	taskClient := body.Request.Ctx.Get(metadataTaskClient)
	taskRetries, _ := strconv.Atoi(body.Request.Ctx.Get(metadataTaskRetries))
	cache.GuildSearch.AddRecord([]string{region, query}, guildProfiles, status, taskRetries, taskClient, startTime)
}
