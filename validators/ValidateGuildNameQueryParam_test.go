package validators

import "testing"

func TestValidateGuildNameQueryParam(t *testing.T) {
	tests := []struct {
		expectedName string
		expectedOk   bool
		input        []string
	}{
		{input: []string{"1NumberGuild"}, expectedName: "1numberguild", expectedOk: true}, // Contains a number
		{input: []string{"Adventure_Guild"}, expectedName: "adventure_guild", expectedOk: true},
		{input: []string{"FirstGuild", "SecondGuild"}, expectedName: "firstguild", expectedOk: true},
		{input: []string{"MyGuild"}, expectedName: "myguild", expectedOk: true},
		{input: []string{"고대신"}, expectedName: "고대신", expectedOk: true}, // Guild name with Korean characters

		{input: []string{""}, expectedName: "", expectedOk: false},                                       // Empty guild name
		{input: []string{"A Guild With Spaces"}, expectedName: "a guild with spaces", expectedOk: false}, // Contains spaces
		{input: []string{"Some$"}, expectedName: "some$", expectedOk: false},                             // Contains an invalid symbol
		{input: []string{"x"}, expectedName: "x", expectedOk: false},                                     // Too short
		{input: []string{}, expectedName: "", expectedOk: false},
	}

	for _, test := range tests {
		name, ok := ValidateGuildNameQueryParam(test.input)
		if name != test.expectedName || ok != test.expectedOk {
			t.Errorf("Input: %v, Expected: %v %v, Got: %v %v", test.input, test.expectedName, test.expectedOk, name, ok)
		}
	}
}
