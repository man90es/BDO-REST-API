package scrapers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/models"
)

func ScrapeAdventurer(profileTarget string) (profile models.Profile, status int) {
	c := collyFactory()

	profile.ProfileTarget = profileTarget
	status = http.StatusNotFound

	c.OnHTML(`.closetime_message`, func(e *colly.HTMLElement) {
		status = http.StatusServiceUnavailable
	})

	c.OnHTML(`.nick`, func(e *colly.HTMLElement) {
		profile.FamilyName = e.Text
		status = http.StatusOK
	})

	c.OnHTML(`.region_info`, func(e *colly.HTMLElement) {
		profile.Region = e.Text
	})

	c.OnHTML(`.desc.guild a`, func(e *colly.HTMLElement) {
		if e.Attr("href") != "javscript:void(0)" {
			profile.Guild = &models.GuildProfile{
				Name: e.Text,
			}
		}
	})

	c.OnHTML(`.desc.guild span`, func(e *colly.HTMLElement) {
		profile.Privacy = profile.Privacy | models.PrivateGuild
	})

	c.OnHTML(`.line_list .desc:not(.guild)`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.Text))
		profile.CreatedOn = &createdOn
	})

	c.OnHTML(`.character_desc_area .character_info span:nth-child(3) em`, func(e *colly.HTMLElement) {
		if e.Text != "Private" {
			contributionPoints, _ := strconv.Atoi(e.Text)
			profile.ContributionPoints = uint16(contributionPoints)
		} else {
			profile.Privacy = profile.Privacy | models.PrivateContrib
		}
	})

	c.OnHTML(`.character_desc_area`, func(e *colly.HTMLElement) {
		character := models.Character{
			Class: e.ChildText(".character_info .character_symbol em:last-child"),
		}

		e.ForEach(`.selected_label`, func(ind int, el *colly.HTMLElement) {
			character.Main = true
		})

		if levelStr := e.ChildText(".character_info span:nth-child(2) em"); levelStr != "Private" {
			level, _ := strconv.Atoi(levelStr)
			character.Level = uint8(level)
		} else {
			profile.Privacy = profile.Privacy | models.PrivateLevel
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
				character.SpecLevels = &models.Specs{
					Gathering:  specLevels[0],
					Fishing:    specLevels[1],
					Hunting:    specLevels[2],
					Cooking:    specLevels[3],
					Alchemy:    specLevels[4],
					Processing: specLevels[5],
					Training:   specLevels[6],
					Trading:    specLevels[7],
					Farming:    specLevels[8],
					Sailing:    specLevels[9],
					Barter:     specLevels[10],
				}
			}
		}

		profile.Characters = append(profile.Characters, character)
	})

	c.OnHTML(`.character_spec.lock`, func(e *colly.HTMLElement) {
		profile.Privacy = profile.Privacy | models.PrivateSpecs
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure/Profile?profileTarget=%v", profileTarget))

	return
}
