package validators

import (
	"testing"
)

func TestValidatePageQueryParam(t *testing.T) {
	tests := []struct {
		input    []string
		expected uint16
	}{
		{input: []string{}, expected: 1},
		{input: []string{"5"}, expected: 5},
		{input: []string{"invalid"}, expected: 1},
		{input: []string{"32767"}, expected: 32767},
		{input: []string{"-7"}, expected: 1},
	}

	for _, test := range tests {
		result := ValidatePageQueryParam(test.input)
		if result != test.expected {
			t.Errorf("For input %v, expected %d, but got %d", test.input, test.expected, result)
		}
	}
}
