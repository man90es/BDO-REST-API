package scrapers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/models"
)

func ScrapeGuild(region, name string) (guildProfile models.GuildProfile, status int) {
	c := collyFactory()

	status = http.StatusNotFound

	c.OnHTML(closetimeSelector, func(e *colly.HTMLElement) {
		status = http.StatusServiceUnavailable
	})

	c.OnHTML(`.region_info`, func(e *colly.HTMLElement) {
		guildProfile.Region = e.Text
	})

	c.OnHTML(`.guild_name p`, func(e *colly.HTMLElement) {
		guildProfile.Name = e.Text
		status = http.StatusOK
	})

	c.OnHTML(`.line_list.mob_none .desc`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.Text))
		guildProfile.CreatedOn = &createdOn
	})

	c.OnHTML(`.line_list:not(.mob_none) li:nth-child(2) .desc .text a`, func(e *colly.HTMLElement) {
		guildProfile.Master = &models.Profile{
			FamilyName:    e.Text,
			ProfileTarget: extractProfileTarget(e.Attr("href")),
		}
	})

	c.OnHTML(`.line_list:not(.mob_none) li:nth-child(3) em`, func(e *colly.HTMLElement) {
		population, _ := strconv.Atoi(e.Text)
		guildProfile.Population = uint8(population)
	})

	c.OnHTML(`.line_list:not(.mob_none) li:last-child .desc`, func(e *colly.HTMLElement) {
		text := dry(e.Text)
		if text != "None" && text != "N/A" {
			guildProfile.Occupying = text
		}
	})

	c.OnHTML(`.box_list_area .adventure_list_table a`, func(e *colly.HTMLElement) {
		member := models.Profile{
			FamilyName:    e.Text,
			ProfileTarget: extractProfileTarget(e.Attr("href")),
		}

		guildProfile.Members = append(guildProfile.Members, member)
	})

	c.Visit(fmt.Sprintf("%v/Adventure/Guild/GuildProfile?guildName=%v&region=%v", getSiteRoot(region), name, region))

	return
}
