package validators

import "testing"

func TestValidateAdventurerNameQueryParam(t *testing.T) {
	tests := []struct {
		expectedName string
		expectedOk   bool
		input        []string
	}{
		{input: []string{"1Number"}, expectedName: "1Number", expectedOk: true}, // Starts with a number
		{input: []string{"Adventurer_123"}, expectedName: "Adventurer_123", expectedOk: true},
		{input: []string{"JohnDoe"}, expectedName: "JohnDoe", expectedOk: true},
		{input: []string{"Name1", "Name2"}, expectedName: "Name1", expectedOk: true},
		{input: []string{"고대신"}, expectedName: "고대신", expectedOk: true}, // Adventurer name with Korean characters

		{input: []string{""}, expectedName: "", expectedOk: false},                                                     // Empty adventurer name
		{input: []string{"Ad"}, expectedName: "Ad", expectedOk: false},                                                 // Too short
		{input: []string{"Adventurer With Spaces"}, expectedName: "Adventurer With Spaces", expectedOk: false},         // Contains spaces
		{input: []string{"AdventurerNameTooLong12345"}, expectedName: "AdventurerNameTooLong12345", expectedOk: false}, // Too long
		{input: []string{"Name$"}, expectedName: "Name$", expectedOk: false},                                           // Contains an invalid symbol
		{input: []string{}, expectedName: "", expectedOk: false},
	}

	for _, test := range tests {
		name, ok := ValidateAdventurerNameQueryParam(test.input)
		if name != test.expectedName || ok != test.expectedOk {
			t.Errorf("Input: %v, Expected: %v %v, Got: %v %v", test.input, test.expectedName, test.expectedOk, name, ok)
		}
	}
}
