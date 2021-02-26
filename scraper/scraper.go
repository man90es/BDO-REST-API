package scraper

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"regexp"

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
		guildProfile.GuildMaster = &entity.Profile{
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

func ScrapeProfile(profileTarget string) (profile entity.Profile)  {
	c := colly.NewCollector()

	profile.ProfileTarget = profileTarget

	c.OnHTML(`.nick`, func(e *colly.HTMLElement) {
		profile.FamilyName = e.Text
	})

	c.OnHTML(`.region_info`, func(e *colly.HTMLElement) {
		profile.Region = e.Text
	})

	c.OnHTML(`.desc.guild a`, func(e *colly.HTMLElement) {
		if e.Attr("href") != "javscript:void(0)" {
			profile.Guild = &entity.GuildProfile{
				Name: e.Text,
				Region: profile.Region,
			}
		}
	})

	c.OnHTML(`.line_list .desc:not(.guild)`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.Text))
		profile.CreatedOn = &createdOn
	})

	c.OnHTML(`.character_desc_area .character_info span:nth-child(3) em`, func(e *colly.HTMLElement) {
		if e.Text != "Private" {
			contributionPoints, _ := strconv.Atoi(e.Text)
			profile.ContributionPoints = int16(contributionPoints)
		}
	})

	c.OnHTML(`.character_desc_area`, func(e *colly.HTMLElement) {
		levelStr := e.ChildText(".character_info span:nth-child(2) em")
		var level int

		if levelStr != "Private" {
			level, _ = strconv.Atoi(levelStr)
		}
		name := e.ChildText(".character_name")
		nameEndIndex := strings.Index(name, "\n")

		if nameEndIndex > -1 {
			name = name[:nameEndIndex]
		}

		character := entity.Character{
			Name: name,
			Class: e.ChildText(".character_info .character_symbol em:last-child"),
			Level: int8(level),
		}

		var specLevels []string

		e.ForEach(".character_spec:not(.lock) .spec_level", func(_ int, el *colly.HTMLElement) {
			// "Beginner1" → "Beginner 1"
			i := regexp.MustCompile(`[0-9]`).FindStringIndex(el.Text)[0]
			level := el.Text[:i] + " " + el.Text[i:]

			specLevels = append(specLevels, level)
		})

		if len(specLevels) > 0 {
			character.SpecLevels = &entity.Specs{
				specLevels[0],
				specLevels[1],
				specLevels[2],
				specLevels[3],
				specLevels[4],
				specLevels[5],
				specLevels[6],
				specLevels[7],
				specLevels[8],
				specLevels[9],
				specLevels[10],
			}
		}


		profile.Characters = append(profile.Characters, character)
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure/Profile?profileTarget=%v", profileTarget))

	return
}
