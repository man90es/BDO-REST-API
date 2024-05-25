package scrapers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/models"
	"bdo-rest-api/translators"
)

func ScrapeAdventurerSearch(region string, query string, searchType uint8, page uint16) (profiles []models.Profile, status int) {
	c := newScraper(region)

	status = http.StatusNotFound

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		status = http.StatusOK
		profile := models.Profile{
			Region:        region,
			FamilyName:    e.ChildText(".title a"),
			ProfileTarget: extractProfileTarget(e.ChildAttr(".title a", "href")),
		}

		if region != "SA" && region != "KR" {
			profile.Region = e.ChildText(".region_info")
		}

		if len(e.ChildAttr(".state a", "href")) > 0 {
			profile.Guild = &models.GuildProfile{
				Name: e.ChildText(".state a"),
			}
		}

		// Sometimes site displays text "You have not set your main character."
		// instead of a character
		if len(e.ChildText(".name")) > 0 {
			profile.Characters = make([]models.Character, 1)

			profile.Characters[0].Class = e.ChildText(".name")
			profile.Characters[0].Name = e.ChildText(".text")

			// Site displays the main character when searching by family name
			// And the searched character when searching by character name
			if searchType == 2 {
				profile.Characters[0].Main = true
			}

			if region != "EU" && region != "NA" {
				translators.TranslateClassName(&profile.Characters[0].Class)
			}

			if level, err := strconv.Atoi(e.ChildText(".level")[3:]); err == nil {
				profile.Characters[0].Level = uint8(level)
			}
		}

		profiles = append(profiles, profile)
	})

	c.Visit(fmt.Sprintf("?region=%v&searchType=%v&searchKeyword=%v&Page=%v", region, searchType, query, page))

	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		status = http.StatusServiceUnavailable
	}

	return
}
