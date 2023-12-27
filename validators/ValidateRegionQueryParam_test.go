package validators

import (
	"testing"
)

func TestValidateRegionQueryParameter(t *testing.T) {
	tests := []struct {
		expectedOk     bool
		expectedRegion string
		input          []string
	}{
		{input: []string{}, expectedRegion: "EU", expectedOk: true},
		{input: []string{"NA"}, expectedRegion: "NA", expectedOk: true},
		{input: []string{"na"}, expectedRegion: "NA", expectedOk: true},
		{input: []string{"SA"}, expectedRegion: "SA", expectedOk: true},
		{input: []string{"EU"}, expectedRegion: "EU", expectedOk: true},
		{input: []string{"KR"}, expectedRegion: "KR", expectedOk: false},
		{input: []string{"NA", "SA"}, expectedRegion: "NA", expectedOk: true}, // Takes the first region in case of multiple regions
	}

	for _, test := range tests {
		result, ok := ValidateRegionQueryParam(test.input)
		if result != test.expectedRegion || ok != test.expectedOk {
			t.Errorf("For input %v, expected %v %v, but got %v %v", test.input, test.expectedRegion, test.expectedOk, result, ok)
		}
	}
}
