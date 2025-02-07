package validators

import "testing"

func TestValidateSearchTypeQueryParam(t *testing.T) {
	tests := []struct {
		expected string
		input    []string
	}{
		{input: []string{}, expected: "2"},
		{input: []string{"characterName"}, expected: "1"},
		{input: []string{"invalidType"}, expected: "2"},
		{input: []string{"familyName"}, expected: "2"},
	}

	for _, test := range tests {
		result := ValidateSearchTypeQueryParam(test.input)
		if result != test.expected {
			t.Errorf("For input %v, expected %v, but got %v", test.input, test.expected, result)
		}
	}
}
