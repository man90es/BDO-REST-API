package entity

import "time"

type Profile struct {
	FamilyName 			string 		`json:"familyName"`
	ProfileTarget 		string 		`json:"profileTarget"`
	Region 				string 		`json:"region"`
	GuildName 			string 		`json:"guildName,omitempty"`
	ContributionPoints 	int16 		`json:"contributionPoints,omitempty"`
	CreatedOn 			*time.Time 	`json:"createdOn,omitempty"`
	Characters 			[]Character `json:"characters,omitempty"`
}
