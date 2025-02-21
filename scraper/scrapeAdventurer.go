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

func scrapeAdventurer(body *colly.HTMLElement, region string) {
	status := http.StatusNotFound
	profile := models.Profile{
		Region: region,
	}

	profile.ProfileTarget = extractProfileTarget(body.Request.URL.String())

	body.ForEachWithBreak(".nick", func(_ int, e *colly.HTMLElement) bool {
		profile.FamilyName = e.Text
		status = http.StatusOK
		return false
	})

	body.ForEachWithBreak(".desc.guild a", func(_ int, e *colly.HTMLElement) bool {
		profile.Guild = &models.GuildProfile{
			Name: e.Text,
		}
		return false
	})

	body.ForEachWithBreak(".line_list li:nth-child(1) .desc span", func(_ int, e *colly.HTMLElement) bool {
		guildStatus := e.Text

		if region != "EU" && region != "NA" {
			translators.TranslateMisc(&guildStatus)
		}

		if guildStatus == "Private" {
			profile.Privacy = profile.Privacy | models.PrivateGuild
		}

		return false
	})

	body.ForEachWithBreak(".line_list li:nth-child(2) .desc", func(_ int, e *colly.HTMLElement) bool {
		createdOn := utils.ParseDate(e.Text)
		profile.CreatedOn = &createdOn
		return false
	})

	body.ForEachWithBreak(".line_list li:nth-child(3) .desc", func(_ int, e *colly.HTMLElement) bool {
		if contributionPoints, err := strconv.Atoi(e.Text); err == nil {
			profile.ContributionPoints = uint16(contributionPoints)
		} else {
			profile.Privacy = profile.Privacy | models.PrivateContrib
		}

		return false
	})

	body.ForEachWithBreak(".character_spec:not(.lock)", func(_ int, e *colly.HTMLElement) bool {
		specLevels := [11]string{}

		e.ForEach(".spec_level", func(ind int, el *colly.HTMLElement) {
			// "Beginner1" → "Beginner 1"
			i := regexp.MustCompile(`[0-9]`).FindStringIndex(el.Text)[0]
			wordLevel := el.Text[:i]

			if region != "EU" && region != "NA" {
				translators.TranslateSpecLevel(&wordLevel)
			}

			specLevels[ind] = wordLevel + " " + el.Text[i:]
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

		if level, err := strconv.Atoi(e.ChildText(".character_info span:nth-child(2) em")); err == nil {
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

		profile.Characters = append(profile.Characters, character)
	})

	body.ForEachWithBreak(".character_spec.lock", func(_ int, _ *colly.HTMLElement) bool {
		profile.Privacy = profile.Privacy | models.PrivateSpecs
		return false
	})

	if profile.Privacy&models.PrivateLevel == 0 {
		profile.CombatFame = utils.CalculateCombatFame(profile.Characters)
	}

	cache.Profiles.AddRecord([]string{region, profile.ProfileTarget}, profile, status)
}
