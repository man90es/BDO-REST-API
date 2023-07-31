package validators

import "testing"

func TestValidateProfileTarget(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid profile targets with lengths >= 150
		{input: repeat("A", 150), expected: true}, // 150-character string
		{input: repeat("A", 200), expected: true}, // 200-character string

		// Invalid profile targets with lengths < 150
		{input: "", expected: false},
		{input: "Short", expected: false},
		{input: repeat("A", 149), expected: false}, // 149-character string
	}

	for _, test := range tests {
		result := ValidateProfileTarget(&test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}

// Helper function to repeat a character n times.
func repeat(s string, n int) string {
	var result string
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
