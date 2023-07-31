package validators

import "testing"

func TestValidatePage(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid pages
		{input: "1", expected: true},
		{input: "10", expected: true},
		{input: "999", expected: true},

		// Invalid pages
		{input: "-5", expected: false},
		{input: "0", expected: false},
		{input: "abc", expected: false},
	}

	for _, test := range tests {
		result := ValidatePage(&test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
