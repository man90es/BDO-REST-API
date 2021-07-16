package scrapers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/models"
)

func ScrapeGuildSearch(region, query string, page int32) (guildProfiles []models.GuildProfile, status int) {
	c := collyFactory()

	c.OnHTML(`.closetime_message`, func(e *colly.HTMLElement) {
		status = http.StatusServiceUnavailable
	})

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.ChildText(".date")))

		guildProfile := models.GuildProfile{
			Name:   e.ChildText(".guild_title a"),
			Region: e.ChildText(".region_info"),
			Kind:   e.ChildText(".tag_label.guild_label"),
			Master: &models.Profile{
				FamilyName:    e.ChildText(".guild_info a"),
				ProfileTarget: extractProfileTarget(e.ChildAttr(".guild_info a", "href")),
			},
			CreatedOn: &createdOn,
		}

		if membersStr := e.ChildText(".member"); true {
			population, _ := strconv.Atoi(membersStr)
			guildProfile.Population = int16(population)
		}

		guildProfiles = append(guildProfiles, guildProfile)
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure/Guild?region=%v&page=%v&searchText=%v", region, page, query))

	if len(guildProfiles) < 1 {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	return
}
