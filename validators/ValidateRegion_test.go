package validators

import "testing"

func TestValidateRegion(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid regions
		{input: "EU", expected: true},
		{input: "NA", expected: true},
		{input: "SA", expected: true},

		// Invalid regions
		{input: "", expected: false},
		{input: "JP", expected: false},
		{input: "KR", expected: false},
	}

	for _, test := range tests {
		result := ValidateRegion(&test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
