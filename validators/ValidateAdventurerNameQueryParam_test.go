package validators

import "testing"

func TestValidateAdventurerNameQueryParam(t *testing.T) {
	tests := []struct {
		expectedName    string
		expectedOk      bool
		expectedMessage string
		input           []string
		region          string
		searchType      string
	}{
		{input: []string{"1Number"}, region: "EU", searchType: "1", expectedName: "1number", expectedOk: true, expectedMessage: ""},
		{input: []string{"Adventurer_123"}, region: "EU", searchType: "1", expectedName: "adventurer_123", expectedOk: true, expectedMessage: ""},
		{input: []string{"JohnDoe"}, region: "EU", searchType: "1", expectedName: "johndoe", expectedOk: true, expectedMessage: ""},
		{input: []string{"Name1", "Name2"}, region: "EU", searchType: "1", expectedName: "name1", expectedOk: true, expectedMessage: ""},
		{input: []string{"고대신"}, region: "EU", searchType: "1", expectedName: "고대신", expectedOk: true, expectedMessage: ""},

		{input: []string{""}, region: "EU", searchType: "1", expectedName: "", expectedOk: false, expectedMessage: "Adventurer name on EU should be at least 3 symbols long"},
		{input: []string{"With Spaces"}, region: "EU", searchType: "1", expectedName: "with spaces", expectedOk: false, expectedMessage: "Adventurer name contains a forbidden symbol at position 5: ' '"},
		{input: []string{"AdventurerNameTooLong12345"}, region: "EU", searchType: "1", expectedName: "adventurernametoolong12345", expectedOk: false, expectedMessage: "Adventurer name on EU should be at most 19 symbols long"},
		{input: []string{"Name$"}, region: "EU", searchType: "1", expectedName: "name$", expectedOk: false, expectedMessage: "Adventurer name contains a forbidden symbol at position 5: '$'"},
		{input: []string{}, region: "EU", searchType: "1", expectedName: "", expectedOk: false, expectedMessage: "Adventurer name is missing from request"},

		{input: []string{""}, region: "SA", searchType: "2", expectedName: "", expectedOk: false, expectedMessage: "Adventurer name on SA should be at least 2 symbols long"},
		{input: []string{"Ad"}, region: "SA", searchType: "2", expectedName: "ad", expectedOk: true, expectedMessage: ""},
	}

	for _, test := range tests {
		name, ok, message := ValidateAdventurerNameQueryParam(test.input, test.region, test.searchType)
		if name != test.expectedName || ok != test.expectedOk || message != test.expectedMessage {
			t.Errorf("Input: %v %v, Expected: %v %v %v, Got: %v %v %v", test.input, test.region, test.expectedName, test.expectedOk, test.expectedMessage, name, ok, message)
		}
	}
}
