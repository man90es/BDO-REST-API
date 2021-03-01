package scraper

import (
	"fmt"
	"time"
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"

	"gitlab.com/man90/black-desert-social-rest-api/entity"
)

func ScrapeGuildProfileSearch(region, query string, page int32) (guildProfiles []entity.GuildProfile)  {
	c := colly.NewCollector()
	c.SetRequestTimeout(time.Minute / 2)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.ChildText(".date")))

		guildProfile := entity.GuildProfile{
			Name: e.ChildText(".guild_title a"),
			Region: region,
			Kind: e.ChildText(".tag_label.guild_label"),
			GuildMaster: &entity.Profile{
				FamilyName: e.ChildText(".box_list_area li .character_desc a"),
				ProfileTarget: e.ChildAttr(".box_list_area li .character_desc a", "href")[nice:],
				Region: region,
			},
			CreatedOn: &createdOn,
		}

		if membersStr := e.ChildText(".guild_member"); true {
			population, _ := strconv.Atoi(membersStr[9:])
			guildProfile.Population = int16(population)
		}

		guildProfiles = append(guildProfiles, guildProfile)
	})


	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure/Guild?region=%v&page=%v&searchText=%v", region, page, query))

	return
}
