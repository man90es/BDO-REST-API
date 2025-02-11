package scrapers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"

	"bdo-rest-api/cache"
	"bdo-rest-api/models"
	"bdo-rest-api/translators"
	"bdo-rest-api/utils"
)

func ScrapeAdventurer(region string, profileTarget string) (profile models.Profile, status int, date string, expires string) {
	c := newScraper(region)

	profile.ProfileTarget = profileTarget
	profile.Region = region
	status = http.StatusNotFound

	c.OnHTML(`.nick`, func(e *colly.HTMLElement) {
		profile.FamilyName = e.Text
		status = http.StatusOK
	})

	c.OnHTML(`.region_info`, func(e *colly.HTMLElement) {
		profile.Region = e.Text
	})

	c.OnHTML(`.desc.guild a`, func(e *colly.HTMLElement) {
		profile.Guild = &models.GuildProfile{
			Name: e.Text,
		}
	})

	c.OnHTML(`.line_list li:nth-child(1) .desc span`, func(e *colly.HTMLElement) {
		guildStatus := e.Text

		if region != "EU" && region != "NA" {
			translators.TranslateMisc(&guildStatus)
		}

		if guildStatus == "Private" {
			profile.Privacy = profile.Privacy | models.PrivateGuild
		}
	})

	c.OnHTML(`.line_list li:nth-child(2) .desc`, func(e *colly.HTMLElement) {
		createdOn := utils.ParseDate(e.Text)
		profile.CreatedOn = &createdOn
	})

	c.OnHTML(`.line_list li:nth-child(3) .desc`, func(e *colly.HTMLElement) {
		if contributionPoints, err := strconv.Atoi(e.Text); err == nil {
			profile.ContributionPoints = uint16(contributionPoints)
		} else {
			profile.Privacy = profile.Privacy | models.PrivateContrib
		}
	})

	c.OnHTML(`.character_spec:not(.lock)`, func(e *colly.HTMLElement) {
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
	})

	c.OnHTML(`.character_desc_area`, func(e *colly.HTMLElement) {
		character := models.Character{
			Class: e.ChildText(".character_info .character_symbol em:last-child"),
		}

		if region != "EU" && region != "NA" {
			translators.TranslateClassName(&character.Class)
		}

		e.ForEach(`.selected_label`, func(ind int, el *colly.HTMLElement) {
			character.Main = true
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

	c.OnHTML(`.character_spec.lock`, func(e *colly.HTMLElement) {
		profile.Privacy = profile.Privacy | models.PrivateSpecs
	})

	c.Visit(fmt.Sprintf("/Profile?profileTarget=%v", url.QueryEscape(profileTarget)))

	if isCloseTime, _ := GetCloseTime(region); isCloseTime {
		status = http.StatusServiceUnavailable
		date, expires = cache.Profiles.SignalMaintenance([]string{region, profileTarget}, profile, status)
		return
	}

	if profile.Privacy&models.PrivateLevel == 0 {
		profile.CombatFame = utils.CalculateCombatFame(profile.Characters)
	}

	date, expires = cache.Profiles.AddRecord([]string{region, profileTarget}, profile, status)
	return
}
