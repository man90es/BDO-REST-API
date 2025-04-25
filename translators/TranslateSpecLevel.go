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
	"장인":           "Artisan",
	"전문":           "Professional",
	"초급":           "Beginner",
	"명장":           "Master",
	"도인":           "Guru",
}

func TranslateSpecLevel(specLevel *string) {
	if val, ok := specLevelTranslationMap[*specLevel]; ok {
		*specLevel = val
	}
}
