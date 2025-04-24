package scraper

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/translators"
	"bdo-rest-api/utils"
)

func scrapeAdventurer(body *colly.HTMLElement, region, profileTarget string) {
	status := http.StatusNotFound
	profile := models.Profile{
		ProfileTarget: profileTarget,
		Region:        region,
	}

	body.ForEachWithBreak(".nick", func(_ int, e *colly.HTMLElement) bool {
		profile.FamilyName = e.Text
		status = http.StatusOK
		return false
	})

	body.ForEachWithBreak(".lock", func(_ int, _ *colly.HTMLElement) bool {
		// FIXME: This is a remains from granular privacy,
		// boolean would be more straightforward now
		profile.Privacy = 15
		return false
	})

	body.ForEachWithBreak(".profile_detail .desc", func(i int, e *colly.HTMLElement) bool {
		switch i {
		case 0:
			createdOn := utils.ParseDate(e.Text)
			profile.CreatedOn = &createdOn
		case 1:
			text := utils.RemoveExtraSpaces(e.Text)
			translators.TranslateMisc(&text)

			if text != "Not in a guild" {
				profile.Guild = &models.GuildProfile{
					Name: text,
				}
			}
		case 2:
			if gs, err := strconv.Atoi(e.Text); err == nil {
				profile.GS = uint16(gs)
			}
		case 3:
			if energy, err := strconv.Atoi(e.Text); err == nil {
				profile.Energy = uint16(energy)
			}
		case 4:
			if contributionPoints, err := strconv.Atoi(e.Text); err == nil {
				profile.ContributionPoints = uint16(contributionPoints)
			}
		}

		return profile.Privacy == 0
	})

	body.ForEachWithBreak(".character_spec", func(_ int, e *colly.HTMLElement) bool {
		specLevels := [11]string{}

		e.ForEach(".spec_level", func(ind int, el *colly.HTMLElement) {
			// "Beginner1" → "Beginner 1"
			lvIndex := regexp.MustCompile("Lv ").FindStringIndex(el.Text)[0]
			wordLevel := el.Text[:lvIndex]

			if region != "EU" && region != "NA" {
				translators.TranslateSpecLevel(&wordLevel)
			}

			specLevels[ind] = wordLevel + el.Text[lvIndex+2:]
		})

		if len(specLevels[0]) > 0 {
			profile.SpecLevels = &models.Specs{
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
			profile.LifeFame = utils.CalculateLifeFame(specLevels)
		}

		return false
	})

	body.ForEach(".character_desc_area", func(_ int, e *colly.HTMLElement) {
		character := models.Character{
			Class: e.ChildText(".character_info .character_symbol em:last-child"),
		}

		if region != "EU" && region != "NA" {
			translators.TranslateClassName(&character.Class)
		}

		e.ForEachWithBreak(".selected_label", func(_ int, _ *colly.HTMLElement) bool {
			character.Main = true
			return false
		})

		if level, err := strconv.Atoi(e.ChildText(".character_info span:nth-child(2) em:not(.lock)")); err == nil {
			character.Level = uint8(level)
		} else {
			// FIXME: This is a remains of times when privacy had granularity
			profile.Privacy = 15
		}

		if name := e.ChildText(".character_name"); true {
			nameEndIndex := strings.Index(name, "\n")

			if nameEndIndex > -1 {
				character.Name = name[:nameEndIndex]
			} else {
				character.Name = name
			}
		}

		profile.Characters = append(profile.Characters, character)
	})

	if profile.Privacy == 0 {
		profile.CombatFame = utils.CalculateCombatFame(profile.Characters)
	}

	cache.Profiles.AddRecord([]string{region, profileTarget}, profile, status, body.Request.Ctx.Get("taskId"))
}
