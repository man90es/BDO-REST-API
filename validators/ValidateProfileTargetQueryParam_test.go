package validators

import "testing"

func TestValidateProfileTargetQueryParam(t *testing.T) {
	tests := []struct {
		expectedOk bool
		expectedPT string
		input      []string
	}{
		// Valid profile targets with lengths >= 150
		{input: []string{repeat("A", 150)}, expectedPT: repeat("A", 150), expectedOk: true},
		{input: []string{repeat("A", 200)}, expectedPT: repeat("A", 200), expectedOk: true},

		// Invalid profile targets with lengths < 150
		{input: []string{""}, expectedPT: "", expectedOk: false},
		{input: []string{"Short"}, expectedPT: "Short", expectedOk: false},
		{input: []string{repeat("A", 149)}, expectedPT: repeat("A", 149), expectedOk: false},

		// Query param not provided
		{input: []string{}, expectedPT: "", expectedOk: false},

		// Several profileTargets provided
		{input: []string{repeat("A", 150), repeat("B", 150)}, expectedPT: repeat("A", 150), expectedOk: true},
	}

	for _, test := range tests {
		pT, ok := ValidateProfileTargetQueryParam(test.input)
		if pT != test.expectedPT || ok != test.expectedOk {
			t.Errorf("Input: %v, Expected: %v %v, Got: %v %v", test.input, test.expectedPT, test.expectedOk, pT, ok)
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
