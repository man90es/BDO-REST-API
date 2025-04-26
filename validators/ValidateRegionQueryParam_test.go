package validators

import (
	"testing"
)

func TestValidateRegionQueryParameter(t *testing.T) {
	tests := []struct {
		expectedOk      bool
		expectedRegion  string
		expectedMessage string
		input           []string
	}{
		{input: []string{}, expectedRegion: "EU", expectedOk: true, expectedMessage: ""},
		{input: []string{"NA"}, expectedRegion: "NA", expectedOk: true, expectedMessage: ""},
		{input: []string{"na"}, expectedRegion: "NA", expectedOk: true, expectedMessage: ""},
		{input: []string{"SA"}, expectedRegion: "SA", expectedOk: true, expectedMessage: ""},
		{input: []string{"EU"}, expectedRegion: "EU", expectedOk: true, expectedMessage: ""},
		{input: []string{"ABC"}, expectedRegion: "ABC", expectedOk: false, expectedMessage: "Region ABC is not supported"},
		{input: []string{"KR", "SA"}, expectedRegion: "KR", expectedOk: true, expectedMessage: ""}, // Takes the first region in case of multiple regions
	}

	for _, test := range tests {
		result, ok, message := ValidateRegionQueryParam(test.input)
		if result != test.expectedRegion || ok != test.expectedOk || message != test.expectedMessage {
			t.Errorf("For input %v, expected %v %v %v, but got %v %v %v", test.input, test.expectedRegion, test.expectedOk, test.expectedMessage, result, ok, message)
		}
	}
}
