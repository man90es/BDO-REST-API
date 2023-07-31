package validators

import "testing"

func TestValidateAdventurerName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid adventurer names
		{input: "1Number", expected: true}, // Starts with a number
		{input: "Adventurer_123", expected: true},
		{input: "JohnDoe", expected: true},
		{input: "고대신", expected: true}, // Adventurer name with Korean characters

		// Invalid adventurer names
		{input: "", expected: false},                           // Empty adventurer name
		{input: "Ad", expected: false},                         // Too short
		{input: "Adventurer With Spaces", expected: false},     // Contains spaces
		{input: "AdventurerNameTooLong12345", expected: false}, // Too long
		{input: "Name$", expected: false},                      // Contains an invalid symbol '$'
	}

	for _, test := range tests {
		result := ValidateAdventurerName(&test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
