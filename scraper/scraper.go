package scraper

import (
	"fmt"
	"time"
	"strings"

	"github.com/gocolly/colly/v2"

	"gitlab.com/man90/black-desert-social-rest-api/entity"
)

const nice = 69

func dry(s string) string {
    return strings.Join(strings.Fields(s), " ")
}

func ScrapeGuildProfile(region, name string) (guildProfile entity.GuildProfile)  {
	c := colly.NewCollector()

	guildProfile.Region = region
	guildProfile.Name = name

	c.OnHTML(`.line_list.mob_none .desc`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.Text))
		guildProfile.CreatedOn = &createdOn
	})

	c.OnHTML(`.line_list:not(.mob_none) li:nth-child(2) .desc .text a`, func(e *colly.HTMLElement) {
		guildProfile.GuildMaster = entity.Profile{
			FamilyName: e.Text,
			ProfileTarget: e.Attr("href")[nice:],
			Region: region,
		}
	})

	c.OnHTML(`.box_list_area .adventure_list_table a`, func(e *colly.HTMLElement) {
		member := entity.Profile{
			FamilyName: e.Text,
			ProfileTarget: e.Attr("href")[nice:],
			Region: region,
		}

		guildProfile.Members = append(guildProfile.Members, member)
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure/Guild/GuildProfile?guildName=%v&region=%v", name, region))

	return
}
