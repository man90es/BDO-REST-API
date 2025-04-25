package scraper

import (
	"fmt"
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
			text := strings.TrimSpace(e.Text)
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

	body.ForEachWithBreak(".history_level", func(i int, e *colly.HTMLElement) bool {
		if profile.Privacy > 0 {
			return false
		}

		switch i {
		case 0:
			profile.History = &models.History{}
			if mobs, err := strconv.ParseUint(e.Text, 10, 32); err == nil {
				profile.History.Mobs = uint(mobs)
			}
		case 1:
			if fish, err := strconv.ParseUint(e.Text, 10, 32); err == nil {
				profile.History.Fish = uint(fish)
			}
		case 2:
			if loot, err := strconv.ParseUint(e.Text, 10, 32); err == nil {
				profile.History.Loot = uint(loot)
			}
		case 3:
			text := strings.Replace(e.Text[0:len(e.Text)-3], ",", ".", 1)

			if lootWeight, err := strconv.ParseFloat(text, 32); err == nil {
				profile.History.LootWeight = float32(lootWeight)
			}
		}

		return true
	})

	body.ForEachWithBreak(".spec_level", func(i int, e *colly.HTMLElement) bool {
		if profile.Privacy > 0 {
			return false
		}

		// "Beginner1" → "Beginner 1"
		lvIndex := regexp.MustCompile("Lv ").FindStringIndex(e.Text)[0]
		wordLevel := e.Text[:lvIndex]

		if region != "EU" && region != "NA" {
			translators.TranslateSpecLevel(&wordLevel)
		}

		value := wordLevel + e.Text[lvIndex+2:]

		switch i {
		case 0:
			profile.SpecLevels = &models.Specs{}
			profile.SpecLevels.Gathering = value
		case 1:
			profile.SpecLevels.Fishing = value
		case 2:
			profile.SpecLevels.Hunting = value
		case 3:
			profile.SpecLevels.Cooking = value
		case 4:
			profile.SpecLevels.Alchemy = value
		case 5:
			profile.SpecLevels.Processing = value
		case 6:
			profile.SpecLevels.Training = value
		case 7:
			profile.SpecLevels.Trading = value
		case 8:
			profile.SpecLevels.Farming = value
		case 9:
			profile.SpecLevels.Sailing = value
		case 10:
			profile.SpecLevels.Barter = value
			profile.LifeFame = utils.CalculateLifeFame(profile.SpecLevels)
		}

		return true
	})

	body.ForEachWithBreak(".spec_stat", func(i int, e *colly.HTMLElement) bool {
		if profile.Privacy > 0 {
			return false
		}

		loot, err := strconv.ParseUint(strings.TrimSpace(e.Text), 10, 16)
		if err != nil {
			fmt.Println(err)
			return false
		}

		value := uint16(loot)

		switch i {
		case 0:
			profile.Mastery = &models.Mastery{}
			profile.Mastery.Gathering = value
		case 1:
			profile.Mastery.Fishing = value
		case 2:
			profile.Mastery.Hunting = value
		case 3:
			profile.Mastery.Cooking = value
		case 4:
			profile.Mastery.Alchemy = value
		case 5:
			profile.Mastery.Processing = value
		case 6:
			profile.Mastery.Training = value
		case 7:
			profile.Mastery.Trading = value
		case 8:
			profile.Mastery.Farming = value
		case 9:
			profile.Mastery.Sailing = value
		case 10:
			profile.Mastery.Barter = value
		}

		return true
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
