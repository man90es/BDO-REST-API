package validators

import "testing"

func TestValidateProfileTargetQueryParam(t *testing.T) {
	tests := []struct {
		expectedOk      bool
		expectedPT      string
		expectedMessage string
		input           []string
	}{
		// Valid profile targets with lengths >= 150
		{input: []string{repeat("A", 150)}, expectedPT: repeat("A", 150), expectedOk: true, expectedMessage: ""},
		{input: []string{repeat("A", 200)}, expectedPT: repeat("A", 200), expectedOk: true, expectedMessage: ""},

		// Invalid profile targets with lengths < 150
		{input: []string{""}, expectedPT: "", expectedOk: false, expectedMessage: "Profile target has to be at least 150 characters long"},
		{input: []string{"Short"}, expectedPT: "Short", expectedOk: false, expectedMessage: "Profile target has to be at least 150 characters long"},
		{input: []string{repeat("A", 149)}, expectedPT: repeat("A", 149), expectedOk: false, expectedMessage: "Profile target has to be at least 150 characters long"},

		// Query param not provided
		{input: []string{}, expectedPT: "", expectedOk: false, expectedMessage: "Profile target is missing from the request"},

		// Several profileTargets provided
		{input: []string{repeat("A", 150), repeat("B", 150)}, expectedPT: repeat("A", 150), expectedOk: true, expectedMessage: ""},
	}

	for _, test := range tests {
		pT, ok, message := ValidateProfileTargetQueryParam(test.input)
		if pT != test.expectedPT || ok != test.expectedOk || message != test.expectedMessage {
			t.Errorf("Input: %v, Expected: %v %v %v, Got: %v %v %v", test.input, test.expectedPT, test.expectedOk, test.expectedMessage, pT, ok, message)
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
