package entity

type Character struct {
	Name 		string 	`json:"name"`
	Class 		string 	`json:"class"`
	Level 		int8 	`json:"level,omitempty"`
	SpecLevels 	*Specs 	`json:"specLevels,omitempty"`
}
