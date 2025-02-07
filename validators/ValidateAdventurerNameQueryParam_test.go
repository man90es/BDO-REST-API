package validators

import "testing"

func TestValidateAdventurerNameQueryParam(t *testing.T) {
	tests := []struct {
		expectedName    string
		expectedOk      bool
		expectedMessage string
		input           []string
		region          string
	}{
		{input: []string{"1Number"}, region: "EU", expectedName: "1number", expectedOk: true, expectedMessage: ""},
		{input: []string{"Adventurer_123"}, region: "EU", expectedName: "adventurer_123", expectedOk: true, expectedMessage: ""},
		{input: []string{"JohnDoe"}, region: "EU", expectedName: "johndoe", expectedOk: true, expectedMessage: ""},
		{input: []string{"Name1", "Name2"}, region: "EU", expectedName: "name1", expectedOk: true, expectedMessage: ""},
		{input: []string{"고대신"}, region: "EU", expectedName: "고대신", expectedOk: true, expectedMessage: ""},

		{input: []string{""}, region: "EU", expectedName: "", expectedOk: false, expectedMessage: "Adventurer name should be between 3 and 16 symbols long"},
		{input: []string{"Ad"}, region: "EU", expectedName: "ad", expectedOk: false, expectedMessage: "Adventurer name should be between 3 and 16 symbols long"},
		{input: []string{"With Spaces"}, region: "EU", expectedName: "with spaces", expectedOk: false, expectedMessage: "Adventurer name contains a forbidden symbol at position 5: ' '"},
		{input: []string{"AdventurerNameTooLong12345"}, region: "EU", expectedName: "adventurernametoolong12345", expectedOk: false, expectedMessage: "Adventurer name should be between 3 and 16 symbols long"},
		{input: []string{"Name$"}, region: "EU", expectedName: "name$", expectedOk: false, expectedMessage: "Adventurer name contains a forbidden symbol at position 5: '$'"},
		{input: []string{}, region: "EU", expectedName: "", expectedOk: false, expectedMessage: "Adventurer name is missing from request"},

		{input: []string{""}, region: "SA", expectedName: "", expectedOk: false, expectedMessage: "Adventurer name should be between 2 and 16 symbols long"},
		{input: []string{"Ad"}, region: "SA", expectedName: "ad", expectedOk: true, expectedMessage: ""},
	}

	for _, test := range tests {
		name, ok, message := ValidateAdventurerNameQueryParam(test.input, test.region)
		if name != test.expectedName || ok != test.expectedOk || message != test.expectedMessage {
			t.Errorf("Input: %v %v, Expected: %v %v %v, Got: %v %v %v", test.input, test.region, test.expectedName, test.expectedOk, test.expectedMessage, name, ok, message)
		}
	}
}
