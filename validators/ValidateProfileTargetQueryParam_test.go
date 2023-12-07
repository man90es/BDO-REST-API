package validators

import "testing"

func TestValidateProfileTargetQueryParam(t *testing.T) {
	tests := []struct {
		expectedGood bool
		expectedPT   string
		input        []string
	}{
		// Valid profile targets with lengths >= 150
		{input: []string{repeat("A", 150)}, expectedPT: repeat("A", 150), expectedGood: true},
		{input: []string{repeat("A", 200)}, expectedPT: repeat("A", 200), expectedGood: true},

		// Invalid profile targets with lengths < 150
		{input: []string{""}, expectedPT: "", expectedGood: false},
		{input: []string{"Short"}, expectedPT: "Short", expectedGood: false},
		{input: []string{repeat("A", 149)}, expectedPT: repeat("A", 149), expectedGood: false},

		// Query param not provided
		{input: []string{}, expectedPT: "", expectedGood: false},

		// Several profileTargets provided
		{input: []string{repeat("A", 150), repeat("B", 150)}, expectedPT: repeat("A", 150), expectedGood: true},
	}

	for _, test := range tests {
		pT, good := ValidateProfileTargetQueryParam(test.input)
		if pT != test.expectedPT || good != test.expectedGood {
			t.Errorf("Input: %v, Expected: %v %v, Got: %v %v", test.input, test.expectedPT, test.expectedGood, pT, good)
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
