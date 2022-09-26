package translators

var guildKindTranslationMap = map[string]string{
	"Clã":    "Clan",
	"Guilda": "Guild",
	"길드":     "Guild",
}

func TranslateGuildKind(kind *string) {
	if val, ok := guildKindTranslationMap[*kind]; ok {
		*kind = val
	}
}
