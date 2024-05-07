package utils

import (
	"bdo-rest-api/models"
)

func CalculateCombatFame(characters []models.Character) (combatFame uint32) {
	for _, character := range characters {
		if character.Level < 56 {
			combatFame += uint32(character.Level)
		} else if character.Level < 60 {
			combatFame += uint32(character.Level) * 2
		} else {
			combatFame += uint32(character.Level) * 5
		}
	}

	return combatFame + 1
}
