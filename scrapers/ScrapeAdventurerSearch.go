package scrapers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/models"
)

func ScrapeAdventurerSearch(region, query string, searchType uint8, page uint16) (profiles []models.Profile, status int) {
	c := collyFactory()
	closetime := false

	c.OnHTML(`.closetime_message`, func(e *colly.HTMLElement) {
		closetime = true
	})

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		profile := models.Profile{
			Region:        e.ChildText(".region_info"),
			FamilyName:    e.ChildText(".title a"),
			ProfileTarget: extractProfileTarget(e.ChildAttr(".title a", "href")),
			Characters:    make([]models.Character, 1),
		}

		if e.ChildAttr(".state a", "href") != "javscript:void(0)" {
			profile.Guild = &models.GuildProfile{
				Name: e.ChildText(".state a"),
			}
		}

		profile.Characters[0].Name = e.ChildText(".text")
		profile.Characters[0].Class = e.ChildText(".name")

		if level, err := strconv.Atoi(e.ChildText(".level")[3:]); err == nil {
			profile.Characters[0].Level = uint8(level)
		}

		profiles = append(profiles, profile)
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure?region=%v&searchType=%v&searchKeyword=%v&Page=%v", region, searchType, query, page))

	if closetime {
		status = http.StatusServiceUnavailable
	} else if len(profiles) < 1 {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	return
}
