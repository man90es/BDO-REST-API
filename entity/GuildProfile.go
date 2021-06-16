package entity

import "time"

type GuildProfile struct {
	Name        string     `json:"name"`
	Region      string     `json:"region"`
	Kind        string     `json:"kind,omitempty"`
	CreatedOn   *time.Time `json:"createdOn,omitempty"`
	GuildMaster *Profile   `json:"guildMaster,omitempty"`
	Population  int16      `json:"population,omitempty"`
	Occupying   string     `json:"occupying,omitempty"`
	Members     []Profile  `json:"members,omitempty"`
}
