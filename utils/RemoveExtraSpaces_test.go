package utils

import "testing"

func TestRemoveExtraSpaces(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "This   is a    test   string.",
			expected: "This is a test string.",
		},
		{
			input:    "  A  string  with   extra   spaces  ",
			expected: "A string with extra spaces",
		},
		{
			input:    "Hello\nWorld\n",
			expected: "Hello World",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		result := RemoveExtraSpaces(test.input)
		if result != test.expected {
			t.Errorf("Input: %v, Expected: %v, Got: %v", test.input, test.expected, result)
		}
	}
}
