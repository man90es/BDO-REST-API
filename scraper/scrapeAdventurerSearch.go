package scraper

import (
	"net/http"
	"strconv"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/translators"
)

func scrapeAdventurerSearch(body *colly.HTMLElement, region, query, searchType string) {
	status := http.StatusNotFound
	profiles := make([]models.Profile, 0)

	body.ForEach(".box_list_area li:not(.no_result)", func(_ int, e *colly.HTMLElement) {
		status = http.StatusOK
		profile := models.Profile{
			FamilyName:    e.ChildText(".title a"),
			ProfileTarget: extractProfileTarget(e.ChildAttr(".title a", "href")),
			Region:        region,
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

			if profile.Region != "EU" && profile.Region != "NA" {
				translators.TranslateClassName(&profile.Characters[0].Class)
			}

			if level, err := strconv.Atoi(e.ChildText(".level")[3:]); err == nil {
				profile.Characters[0].Level = uint8(level)
			}
		}

		profiles = append(profiles, profile)
	})

	cache.ProfileSearch.AddRecord([]string{region, query, searchType}, profiles, status)
}
