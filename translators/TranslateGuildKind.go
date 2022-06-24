package translators

var guildKindTranslationMap = map[string]string{
	"Clã":    "Clan",
	"Guilda": "Guild",
}

func TranslateGuildKind(kind *string) {
	if val, ok := guildKindTranslationMap[*kind]; ok {
		*kind = val
	}
}
