package validators

import "testing"

func TestValidateGuildNameQueryParam(t *testing.T) {
	tests := []struct {
		expectedName    string
		expectedOk      bool
		expectedMessage string
		input           []string
		region          string
	}{
		{input: []string{"1NumberGuild"}, region: "EU", expectedName: "1numberguild", expectedOk: true, expectedMessage: ""}, // Contains a number
		{input: []string{"Adventure_Guild"}, region: "NA", expectedName: "adventure_guild", expectedOk: true, expectedMessage: ""},
		{input: []string{"FirstGuild", "SecondGuild"}, region: "EU", expectedName: "firstguild", expectedOk: true, expectedMessage: ""},
		{input: []string{"MyGuild"}, region: "NA", expectedName: "myguild", expectedOk: true, expectedMessage: ""},
		{input: []string{"고대신"}, region: "KR", expectedName: "고대신", expectedOk: true, expectedMessage: ""}, // Guild name with Korean characters
		{input: []string{"IX"}, region: "SA", expectedName: "ix", expectedOk: true, expectedMessage: ""},   // Guild names on SA can be 2 symbols long

		{input: []string{""}, region: "EU", expectedName: "", expectedOk: false, expectedMessage: "Guild name in EU region can't be shorter than 3 symbols"},
		{input: []string{"X"}, region: "SA", expectedName: "x", expectedOk: false, expectedMessage: "Guild name in SA region can't be shorter than 2 symbols"}, // Guild names on SA can be 2 symbols long
		{input: []string{"With Spaces"}, region: "NA", expectedName: "with spaces", expectedOk: false, expectedMessage: "Guild name contains a forbidden symbol at position 5: ' '"},
		{input: []string{"Some$"}, region: "SA", expectedName: "some$", expectedOk: false, expectedMessage: "Guild name contains a forbidden symbol at position 5: '$'"},
		{input: []string{"x"}, region: "KR", expectedName: "x", expectedOk: false, expectedMessage: "Guild name in KR region can't be shorter than 3 symbols"},
		{input: []string{}, region: "EU", expectedName: "", expectedOk: false, expectedMessage: "Guild name is missing from request"},
		{input: []string{"GuildNameThatIsWayTooLong"}, region: "NA", expectedName: "guildnamethatiswaytoolong", expectedOk: false, expectedMessage: "Guild name can't be longer than 16 symbols"},
	}

	for _, test := range tests {
		name, ok, message := ValidateGuildNameQueryParam(test.input, test.region)
		if name != test.expectedName || ok != test.expectedOk || message != test.expectedMessage {
			t.Errorf("Input: %v, Expected: %v %v %v, Got: %v %v %v", test.input, test.expectedName, test.expectedOk, test.expectedMessage, name, ok, message)
		}
	}
}
