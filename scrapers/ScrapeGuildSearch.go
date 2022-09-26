package scrapers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/models"
	"bdo-rest-api/translators"
)

func ScrapeGuildSearch(region, query string, page uint16) (guildProfiles []models.GuildProfile, status int) {
	c := collyFactory()
	closetime := false

	c.OnHTML(closetimeSelector, func(e *colly.HTMLElement) {
		closetime = true
	})

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.ChildText(".date")))

		guildProfile := models.GuildProfile{
			Name:   e.ChildText(".guild_title a"),
			Region: region,
			Kind:   e.ChildText(".tag_label.guild_label"),
			Master: &models.Profile{
				FamilyName:    e.ChildText(".guild_info a"),
				ProfileTarget: extractProfileTarget(e.ChildAttr(".guild_info a", "href")),
			},
			CreatedOn: &createdOn,
		}

		if region != "SA" && region != "KR" {
			guildProfile.Region = e.ChildText(".region_info")
		}

		if region != "EU" && region != "NA" {
			translators.TranslateGuildKind(&guildProfile.Kind)
		}

		if membersStr := e.ChildText(".member"); true {
			population, _ := strconv.Atoi(membersStr)
			guildProfile.Population = uint8(population)
		}

		guildProfiles = append(guildProfiles, guildProfile)
	})

	c.Visit(fmt.Sprintf("%v/Adventure/Guild?region=%v&page=%v&searchText=%v", getSiteRoot(region), region, page, query))

	if closetime {
		status = http.StatusServiceUnavailable
	} else if len(guildProfiles) < 1 {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	return
}
