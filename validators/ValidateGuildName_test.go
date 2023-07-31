package validators

import "testing"

func TestValidateGuildName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid guild names
		{input: "1NumberGuild", expected: true}, // Contains a number
		{input: "Adventure_Guild", expected: true},
		{input: "MyGuild", expected: true},
		{input: "고대신", expected: true}, // Guild name with Korean characters

		// Invalid guild names
		{input: "", expected: false},                    // Empty guild name
		{input: "A Guild With Spaces", expected: false}, // Contains spaces
		{input: "Some$", expected: false},               // Contains an invalid symbol '$'
		{input: "x", expected: false},                   // Too short
	}

	for _, test := range tests {
		result := ValidateGuildName(&test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
