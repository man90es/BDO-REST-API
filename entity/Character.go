package entity

type Character struct {
	Name 		string 				`json:"name"`
	Class 		string 				`json:"class"`
	Level 		int8 				`json:"level"`
	LifeSkills 	map[string]string 	`json:"lifeSkills"`
}
