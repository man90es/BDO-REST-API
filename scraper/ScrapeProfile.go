package scraper

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"regexp"
	"log"

	"github.com/gocolly/colly/v2"

	"gitlab.com/man90/black-desert-social-rest-api/entity"
)

func ScrapeProfile(profileTarget string) (profile entity.Profile)  {
	c := colly.NewCollector()
	c.SetRequestTimeout(time.Minute / 2)

	profile.ProfileTarget = profileTarget

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

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
		character := entity.Character{
			Class: e.ChildText(".character_info .character_symbol em:last-child"),
		}

		if levelStr := e.ChildText(".character_info span:nth-child(2) em"); levelStr != "Private" {
			level, _ := strconv.Atoi(levelStr)
			character.Level = int8(level)
		}

		if name := e.ChildText(".character_name"); true {
			nameEndIndex := strings.Index(name, "\n")

			if nameEndIndex > -1 {
				character.Name = name[:nameEndIndex]
			} else {
				character.Name = name
			}
		} 

		if specLevels := [11]string{}; true {
			e.ForEach(".character_spec:not(.lock) .spec_level", func(ind int, el *colly.HTMLElement) {
				// "Beginner1" → "Beginner 1"
				i := regexp.MustCompile(`[0-9]`).FindStringIndex(el.Text)[0]
				level := el.Text[:i] + " " + el.Text[i:]

				specLevels[ind] = level
			})

			if len(specLevels[0]) > 0 {
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
		}

		profile.Characters = append(profile.Characters, character)
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure/Profile?profileTarget=%v", profileTarget))

	return
}
