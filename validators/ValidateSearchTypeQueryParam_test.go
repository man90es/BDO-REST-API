package validators

import "testing"

func TestValidateSearchTypeQueryParam(t *testing.T) {
	tests := []struct {
		expectedNum uint8
		expectedStr string
		input       []string
	}{
		{input: []string{}, expectedNum: 2, expectedStr: "familyName"},
		{input: []string{"characterName"}, expectedNum: 1, expectedStr: "characterName"},
		{input: []string{"invalidType"}, expectedNum: 2, expectedStr: "familyName"},
		{input: []string{"familyName"}, expectedNum: 2, expectedStr: "familyName"},
	}

	for _, test := range tests {
		resultNum, resultStr := ValidateSearchTypeQueryParam(test.input)
		if resultNum != test.expectedNum || resultStr != test.expectedStr {
			t.Errorf("For input %v, expected %v %v, but got %v %v", test.input, test.expectedNum, test.expectedStr, resultNum, resultStr)
		}
	}
}
