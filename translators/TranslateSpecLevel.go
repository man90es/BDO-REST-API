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
}

func TranslateSpecLevel(specLevel *string) {
	if val, ok := specLevelTranslationMap[*specLevel]; ok {
		*specLevel = val
	}
}
