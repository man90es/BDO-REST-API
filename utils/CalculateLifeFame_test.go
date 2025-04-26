package utils

import (
	"bdo-rest-api/models"
	"testing"
)

func TestCalculateLifeFame(t *testing.T) {
	tests := []struct {
		input    models.Specs
		expected uint16
	}{
		{
			input:    models.Specs{"Artisan 6", "Professional 4", "Artisan 1", "Master 6", "Professional 7", "Artisan 8", "Professional 10", "Apprentice 9", "Skilled 4", "Apprentice 4", "Beginner 1"},
			expected: 907,
		},
		{
			input:    models.Specs{"Guru 7", "Skilled 5", "Beginner 8", "Guru 52", "Guru 27", "Guru 35", "Artisan 3", "Apprentice 7", "Guru 15", "Geginner 6", "Beginner 1"},
			expected: 1738,
		},
		{
			input:    models.Specs{"Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1", "Beginner 1"},
			expected: 1,
		},
	}

	for _, test := range tests {
		result := CalculateLifeFame(&test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
