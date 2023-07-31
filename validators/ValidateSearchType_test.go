package validators

import "testing"

func TestValidateSearchType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid search types
		{input: "characterName", expected: true},
		{input: "familyName", expected: true},

		// Invalid search types
		{input: "", expected: false},
		{input: "invalidType", expected: false},
		{input: "someOtherType", expected: false},
	}

	for _, test := range tests {
		result := ValidateSearchType(&test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
