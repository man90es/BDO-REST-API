package scraper

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/entity"
)

func ScrapeProfileSearch(region, query string, searchType int8, page int32) (profiles []entity.Profile, err scrapedError) {
	c := collyFactory()

	c.OnHTML(`.closetime_message`, func(e *colly.HTMLElement) {
		err = &entity.MaintenanceError{}
	})

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		profile := entity.Profile{
			Region:        e.ChildText(".region_info"),
			FamilyName:    e.ChildText(".title a"),
			ProfileTarget: extractProfileTarget(e.ChildAttr(".title a", "href")),
			Characters:    make([]entity.Character, 1),
		}

		if e.ChildAttr(".state a", "href") != "javscript:void(0)" {
			profile.Guild = &entity.GuildProfile{
				Name: e.ChildText(".state a"),
			}
		}

		profile.Characters[0].Name = e.ChildText(".text")
		profile.Characters[0].Class = e.ChildText(".name")

		if level, err := strconv.Atoi(e.ChildText(".level")[3:]); err == nil {
			profile.Characters[0].Level = int8(level)
		}

		profiles = append(profiles, profile)
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure?region=%v&searchType=%v&searchKeyword=%v&Page=%v", region, searchType, query, page))

	if len(profiles) < 1 {
		err = &entity.NotFoundError{}
	}

	return
}
