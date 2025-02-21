package scraper

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/utils"
)

func scrapeGuild(body *colly.HTMLElement, region string) {
	status := http.StatusNotFound
	guildProfile := models.GuildProfile{
		Region: region,
	}

	body.ForEachWithBreak(".guild_name p", func(_ int, e *colly.HTMLElement) bool {
		guildProfile.Name = e.Text
		status = http.StatusOK
		return false
	})

	body.ForEachWithBreak(".line_list.mob_none .desc", func(_ int, e *colly.HTMLElement) bool {
		createdOn := utils.ParseDate(e.Text)
		guildProfile.CreatedOn = &createdOn
		return false
	})

	body.ForEachWithBreak(".line_list:not(.mob_none) li:nth-child(2) .desc .text a", func(_ int, e *colly.HTMLElement) bool {
		guildProfile.Master = &models.Profile{
			FamilyName:    e.Text,
			ProfileTarget: extractProfileTarget(e.Attr("href")),
		}
		return false
	})

	body.ForEachWithBreak(".line_list:not(.mob_none) li:nth-child(3) em", func(_ int, e *colly.HTMLElement) bool {
		population, _ := strconv.Atoi(e.Text)
		guildProfile.Population = uint8(population)
		return false
	})

	body.ForEachWithBreak(".line_list:not(.mob_none) li:last-child .desc", func(_ int, e *colly.HTMLElement) bool {
		text := utils.RemoveExtraSpaces(e.Text)
		if text != "None" && text != "N/A" && text != "없음" {
			guildProfile.Occupying = text
		}
		return false
	})

	body.ForEach(".box_list_area .adventure_list_table a", func(_ int, e *colly.HTMLElement) {
		member := models.Profile{
			FamilyName:    e.Text,
			ProfileTarget: extractProfileTarget(e.Attr("href")),
		}

		guildProfile.Members = append(guildProfile.Members, member)
	})

	cache.GuildProfiles.AddRecord([]string{region, strings.ToLower(guildProfile.Name)}, guildProfile, status)
}
