package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		// Valid date formats
		{
			input:    "2023.08.01",
			expected: time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "01/08/2023",
			expected: time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "2023-08-01",
			expected: time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "Aug 1, 2023",
			expected: time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
		},

		// Invalid date formats
		{
			input:    "2023.13.31", // Invalid month
			expected: time.Time{},
		},
		{
			input:    "07-31/2023", // Invalid separator
			expected: time.Time{},
		},

		// Empty input
		{
			input:    "",
			expected: time.Time{},
		},
	}

	for _, test := range tests {
		result := ParseDate(test.input)
		if !result.Equal(test.expected) {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
