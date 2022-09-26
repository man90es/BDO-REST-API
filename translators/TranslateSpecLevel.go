package translators

var specLevelTranslationMap = map[string]string{
	"Aprendiz":     "Apprentice",
	"Artesão":      "Artisan",
	"Hábil":        "Skilled",
	"Iniciante":    "Beginner",
	"Mestre":       "Master",
	"Novato":       "Beginner",
	"Proficiente":  "Skilled",
	"Profissional": "Professional",
	"견습":           "Apprentice",
	"숙련":           "Skilled",
	"초급":           "Beginner",
}

func TranslateSpecLevel(specLevel *string) {
	if val, ok := specLevelTranslationMap[*specLevel]; ok {
		*specLevel = val
	}
}
