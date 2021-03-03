package scraper

import (
	"fmt"
	"strconv"
	"time"
	"log"

	"github.com/gocolly/colly/v2"

	"gitlab.com/man90/black-desert-social-rest-api/entity"
)

func ScrapeProfileSearch(region, query string, searchType int8, page int32) (profiles []entity.Profile, err error)  {
	c := colly.NewCollector()
	c.SetRequestTimeout(time.Minute / 2)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(`.closetime_message`, func(e *colly.HTMLElement) {
		err = fmt.Errorf(closetimeMessage)
	})

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		profile := entity.Profile{
			Region: region,
			FamilyName: e.ChildText(".title a"),
			ProfileTarget: e.ChildAttr(".title a", "href")[nice:],
			Characters: make([]entity.Character, 1),
		}

		if e.ChildAttr(".state a", "href") != "javscript:void(0)" {
			profile.Guild = &entity.GuildProfile{
				Name: e.ChildText(".state a"),
				Region: region,
			}
		}

		profile.Characters[0].Name = e.ChildText(".text")
		profile.Characters[0].Class = e.ChildText(".name")

		if level, err := strconv.Atoi(e.ChildText(".level")[3:]); err == nil {
			profile.Characters[0].Level = int8(level)
		}

		profiles = append(profiles, profile)
	})

	if (len(query) < 1) {
		c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure?region=%v&searchType=3&Page=%v", region, page))
	} else {
		c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure?region=%v&searchType=%v&searchKeyword=%v&Page=%v", region, searchType, query, page))
	}

	return
}
