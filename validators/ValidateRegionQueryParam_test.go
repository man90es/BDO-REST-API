package validators

import (
	"testing"
)

func TestValidateRegionQueryParameter(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{input: []string{}, expected: "EU"},
		{input: []string{"NA"}, expected: "NA"},
		{input: []string{"na"}, expected: "NA"},
		{input: []string{"SA"}, expected: "SA"},
		{input: []string{"EU"}, expected: "EU"},
		{input: []string{"KR"}, expected: "EU"},       // Assuming "KR" falls back to "EU" until translations are ready
		{input: []string{"NA", "SA"}, expected: "NA"}, // Takes the first region in case of multiple regions
	}

	for _, test := range tests {
		result := ValidateRegionQueryParam(test.input)
		if result != test.expected {
			t.Errorf("For input %v, expected %s, but got %s", test.input, test.expected, result)
		}
	}
}
