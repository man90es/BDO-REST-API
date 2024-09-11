package validators

import "testing"

func TestValidateAdventurerNameQueryParam(t *testing.T) {
	tests := []struct {
		expectedName string
		expectedOk   bool
		input        []string
		region       string
	}{
		{input: []string{"1Number"}, region: "EU", expectedName: "1Number", expectedOk: true}, // Starts with a number
		{input: []string{"Adventurer_123"}, region: "EU", expectedName: "Adventurer_123", expectedOk: true},
		{input: []string{"JohnDoe"}, region: "EU", expectedName: "JohnDoe", expectedOk: true},
		{input: []string{"Name1", "Name2"}, region: "EU", expectedName: "Name1", expectedOk: true},
		{input: []string{"고대신"}, region: "EU", expectedName: "고대신", expectedOk: true}, // Adventurer name with Korean characters

		{input: []string{""}, region: "EU", expectedName: "", expectedOk: false},                                                     // Empty adventurer name
		{input: []string{"Ad"}, region: "EU", expectedName: "Ad", expectedOk: false},                                                 // Too short
		{input: []string{"Adventurer With Spaces"}, region: "EU", expectedName: "Adventurer With Spaces", expectedOk: false},         // Contains spaces
		{input: []string{"AdventurerNameTooLong12345"}, region: "EU", expectedName: "AdventurerNameTooLong12345", expectedOk: false}, // Too long
		{input: []string{"Name$"}, region: "EU", expectedName: "Name$", expectedOk: false},                                           // Contains an invalid symbol
		{input: []string{}, region: "EU", expectedName: "", expectedOk: false},

		{input: []string{""}, region: "SA", expectedName: "", expectedOk: false},
		{input: []string{"Ad"}, region: "SA", expectedName: "Ad", expectedOk: true},
	}

	for _, test := range tests {
		name, ok := ValidateAdventurerNameQueryParam(test.input, test.region)
		if name != test.expectedName || ok != test.expectedOk {
			t.Errorf("Input: %v %v, Expected: %v %v, Got: %v %v", test.input, test.region, test.expectedName, test.expectedOk, name, ok)
		}
	}
}
