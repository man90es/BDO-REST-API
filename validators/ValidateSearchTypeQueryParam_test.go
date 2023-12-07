package validators

import "testing"

func TestValidateSearchTypeQueryParam(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{input: []string{}, expected: "familyName"},
		{input: []string{"characterName"}, expected: "characterName"},
		{input: []string{"invalidType"}, expected: "familyName"},
		{input: []string{"familyName"}, expected: "familyName"},
	}

	for _, test := range tests {
		result := ValidateSearchTypeQueryParam(test.input)
		if result != test.expected {
			t.Errorf("For input %v, expected %s, but got %s", test.input, test.expected, result)
		}
	}
}
