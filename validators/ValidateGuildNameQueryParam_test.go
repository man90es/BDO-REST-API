package validators

import "testing"

func TestValidateGuildNameQueryParam(t *testing.T) {
	tests := []struct {
		expectedName    string
		expectedOk      bool
		expectedMessage string
		input           []string
	}{
		{input: []string{"1NumberGuild"}, expectedName: "1numberguild", expectedOk: true, expectedMessage: ""}, // Contains a number
		{input: []string{"Adventure_Guild"}, expectedName: "adventure_guild", expectedOk: true, expectedMessage: ""},
		{input: []string{"FirstGuild", "SecondGuild"}, expectedName: "firstguild", expectedOk: true, expectedMessage: ""},
		{input: []string{"MyGuild"}, expectedName: "myguild", expectedOk: true, expectedMessage: ""},
		{input: []string{"고대신"}, expectedName: "고대신", expectedOk: true, expectedMessage: ""}, // Guild name with Korean characters

		{input: []string{""}, expectedName: "", expectedOk: false, expectedMessage: "Guild name can't be shorter than 3 symbols"},
		{input: []string{"With Spaces"}, expectedName: "with spaces", expectedOk: false, expectedMessage: "Guild name contains a forbidden symbol at position 5: ' '"},
		{input: []string{"Some$"}, expectedName: "some$", expectedOk: false, expectedMessage: "Guild name contains a forbidden symbol at position 5: '$'"},
		{input: []string{"x"}, expectedName: "x", expectedOk: false, expectedMessage: "Guild name can't be shorter than 3 symbols"},
		{input: []string{}, expectedName: "", expectedOk: false, expectedMessage: "Guild name is missing from request"},
		{input: []string{"GuildNameThatIsWayTooLong"}, expectedName: "guildnamethatiswaytoolong", expectedOk: false, expectedMessage: "Guild name can't be longer than 16 symbols"},
	}

	for _, test := range tests {
		name, ok, message := ValidateGuildNameQueryParam(test.input)
		if name != test.expectedName || ok != test.expectedOk || message != test.expectedMessage {
			t.Errorf("Input: %v, Expected: %v %v %v, Got: %v %v %v", test.input, test.expectedName, test.expectedOk, test.expectedMessage, name, ok, message)
		}
	}
}
