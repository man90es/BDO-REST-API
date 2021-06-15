package scraper

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"

	"black-desert-social-rest-api/entity"
)

func ScrapeGuildProfileSearch(region, query string, page int32) (guildProfiles []entity.GuildProfile, err error) {
	c := collyFactory()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(`.closetime_message`, func(e *colly.HTMLElement) {
		err = fmt.Errorf(closetimeMessage)
	})

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.ChildText(".date")))

		guildProfile := entity.GuildProfile{
			Name:   e.ChildText(".guild_title a"),
			Region: e.ChildText(".region_info"),
			Kind:   e.ChildText(".tag_label.guild_label"),
			GuildMaster: &entity.Profile{
				FamilyName:    e.ChildText(".box_list_area li .character_desc a"),
				ProfileTarget: e.ChildAttr(".box_list_area li .character_desc a", "href")[nice:],
				Region:        e.ChildText(".region_info"),
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
