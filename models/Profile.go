package models

import "time"

type Profile struct {
	Characters         []Character   `json:"characters,omitempty"`
	CombatFame         uint32        `json:"combatFame,omitempty"`
	ContributionPoints uint16        `json:"contributionPoints,omitempty"`
	CreatedOn          *time.Time    `json:"createdOn,omitempty"`
	Energy             uint16        `json:"energy,omitempty"`
	FamilyName         string        `json:"familyName"`
	GS                 uint16        `json:"gs,omitempty"`
	Guild              *GuildProfile `json:"guild,omitempty"`
	History            *History      `json:"history,omitempty"`
	LifeFame           uint16        `json:"lifeFame,omitempty"`
	Privacy            int8          `json:"privacy,omitempty"`
	ProfileTarget      string        `json:"profileTarget"`
	Region             string        `json:"region,omitempty"`
	SpecLevels         *Specs        `json:"specLevels,omitempty"`
}
