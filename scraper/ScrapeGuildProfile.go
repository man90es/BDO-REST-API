package scraper

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"

	"gitlab.com/man90/black-desert-social-rest-api/entity"
)

func ScrapeGuildProfile(region, name string) (guildProfile entity.GuildProfile, err error) {
	c := collyFactory()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnHTML(`.closetime_message`, func(e *colly.HTMLElement) {
		err = fmt.Errorf(closetimeMessage)
	})

	c.OnHTML(`.region_info`, func(e *colly.HTMLElement) {
		guildProfile.Region = e.Text
	})

	c.OnHTML(`.guild_name p`, func(e *colly.HTMLElement) {
		guildProfile.Name = e.Text
	})

	c.OnHTML(`.line_list.mob_none .desc`, func(e *colly.HTMLElement) {
		createdOn, _ := time.Parse("2006-01-02", dry(e.Text))
		guildProfile.CreatedOn = &createdOn
	})

	c.OnHTML(`.line_list:not(.mob_none) li:nth-child(2) .desc .text a`, func(e *colly.HTMLElement) {
		guildProfile.GuildMaster = &entity.Profile{
			FamilyName:    e.Text,
			ProfileTarget: e.Attr("href")[nice:],
			Region:        region,
		}
	})

	c.OnHTML(`.line_list:not(.mob_none) li:nth-child(3) em`, func(e *colly.HTMLElement) {
		population, _ := strconv.Atoi(e.Text)
		guildProfile.Population = int16(population)
	})

	c.OnHTML(`.line_list:not(.mob_none) li:last-child .desc`, func(e *colly.HTMLElement) {
		text := dry(e.Text)
		if text != "None" {
			guildProfile.Occupying = text
		}
	})

	c.OnHTML(`.box_list_area .adventure_list_table a`, func(e *colly.HTMLElement) {
		member := entity.Profile{
			FamilyName:    e.Text,
			ProfileTarget: e.Attr("href")[nice:],
			Region:        region,
		}

		guildProfile.Members = append(guildProfile.Members, member)
	})

	c.Visit(fmt.Sprintf("https://www.naeu.playblackdesert.com/en-US/Adventure/Guild/GuildProfile?guildName=%v&region=%v", name, region))

	return
}
