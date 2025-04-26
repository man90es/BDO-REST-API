package models

type Mastery struct {
	Gathering  uint16 `json:"gathering"`
	Fishing    uint16 `json:"fishing"`
	Hunting    uint16 `json:"hunting"`
	Cooking    uint16 `json:"cooking"`
	Alchemy    uint16 `json:"alchemy"`
	Processing uint16 `json:"processing"`
	Training   uint16 `json:"training"`
	Trading    uint16 `json:"trading"`
	Farming    uint16 `json:"farming"`
	Sailing    uint16 `json:"sailing"`
	Barter     uint16 `json:"barter"`
}
