package translators

var miscTranslationMap = map[string]string{
	"Não está alistado em nenhuma guilda.": "Not in a guild",
	"가입된 길드가 없습니다.":                        "Not in a guild",
}

func TranslateMisc(guildMembershipStatus *string) {
	if val, ok := miscTranslationMap[*guildMembershipStatus]; ok {
		*guildMembershipStatus = val
	}
}
