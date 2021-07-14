package entity

import "time"

const (
	PrivateLevel   = 1
	PrivateGuild   = 2
	PrivateContrib = 4
	PrivateSpecs   = 8
)

type Profile struct {
	FamilyName         string        `json:"familyName"`
	ProfileTarget      string        `json:"profileTarget"`
	Region             string        `json:"region,omitempty"`
	Guild              *GuildProfile `json:"guild,omitempty"`
	ContributionPoints int16         `json:"contributionPoints,omitempty"`
	CreatedOn          *time.Time    `json:"createdOn,omitempty"`
	Characters         []Character   `json:"characters,omitempty"`
	Privacy            int8          `json:"privacy,omitempty"`
}
