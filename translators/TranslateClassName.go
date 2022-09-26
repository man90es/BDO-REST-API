package translators

var classNameTranslationMap = map[string]string{
	"Arqueiro":             "Archer",
	"Bruxa":                "Witch",
	"Caçadora":             "Ranger",
	"Cavaleira das Trevas": "Dark Knight",
	"Cavaleira Negra":      "Dark Knight",
	"Corsária":             "Corsair",
	"Domadora":             "Tamer",
	"Feiticeira":           "Sorceress",
	"Guardiã":              "Guardian",
	"Guerreiro":            "Warrior",
	"Lutador":              "Striker",
	"Maga":                 "Witch",
	"Mago":                 "Wizard",
	"Mística":              "Mystic",
	"Musah":                "Musa",
	"Sábio":                "Sage",
	"Sagitário":            "Archer",
	"Valquíria":            "Valkyrie",
	"가디언":                  "Guardian",
	"격투가":                  "Striker",
	"금수랑":                  "Tamer",
	"노바":                   "Nova",
	"닌자":                   "Ninja",
	"다크나이트":                "Dark Knight",
	"드라카니아":                "Drakania",
	"란":                    "Lahn",
	"레인저":                  "Ranger",
	"매화":                   "Maehwa",
	"무사":                   "Musa",
	"미스틱":                  "Mystic",
	"발키리":                  "Valkyrie",
	"샤이":                   "Shai",
	"세이지":                  "Sage",
	"소서러":                  "Sorceress",
	"아처":                   "Archer",
	"워리어":                  "Warrior",
	"위자드":                  "Wizard",
	"위치":                   "Witch",
	"자이언트":                 "Berserker",
	"커세어":                  "Corsair",
	"쿠노이치":                 "Kunoichi",
	"하사신":                  "Hashashin",
}

func TranslateClassName(className *string) {
	if val, ok := classNameTranslationMap[*className]; ok {
		*className = val
	}
}
