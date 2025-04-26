package models

type History struct {
	Fish       uint    `json:"fish"`
	Loot       uint    `json:"loot"`
	LootWeight float32 `json:"lootWeight"`
	Mobs       uint    `json:"mobs"`
}
