package models

import "time"

type GuildProfile struct {
	Name       string     `json:"name"`
	Region     string     `json:"region,omitempty"`
	CreatedOn  *time.Time `json:"createdOn,omitempty"`
	Master     *Profile   `json:"master,omitempty"`
	Population uint8      `json:"population,omitempty"`
	Occupying  string     `json:"occupying,omitempty"`
	Members    []Profile  `json:"members,omitempty"`
}
