package scrapers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/translators"
	"bdo-rest-api/utils"
)

func ScrapeGuildSearch(region, query string) (guildProfiles []models.GuildProfile, status int, date string, expires string) {
	c := newScraper(region)

	status = http.StatusNotFound

	c.OnHTML(`.box_list_area li:not(.no_result)`, func(e *colly.HTMLElement) {
		createdOn := utils.ParseDate(e.ChildText(".date"))
		status = http.StatusOK

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

	c.Visit(fmt.Sprintf("/Guild?region=%v&page=1&searchText=%v", region, query))

	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		status = http.StatusServiceUnavailable
		return
	}

	date, expires = cache.GuildSearch.AddRecord([]string{region, query}, guildProfiles, status)
	return
}
