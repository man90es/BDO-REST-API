package models

type Character struct {
	Name       string `json:"name"`
	Class      string `json:"class"`
	Main       bool   `json:"main,omitempty"`
	Level      uint8  `json:"level,omitempty"`
	SpecLevels *Specs `json:"specLevels,omitempty"`
}
