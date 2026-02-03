package avature_test

import (
	"testing"

	"github.com/dhaef/job-scraper/internal/avature"
)

func TestFindElementsByClasses(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "should find bloomberg",
			input:  "https://bloomberg.avature.net/careers",
			output: "bloomberg",
		},
		{
			name:   "should find macquarie",
			input:  "https://recruitment.macquarie.com/en_US/careers",
			output: "macquarie",
		},
		{
			name:   "should find ea",
			input:  "https://jobs.ea.com/en_US/careers",
			output: "ea",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name := avature.ExtractUniqueHostPart(tc.input)

			if name != tc.output {
				t.Fatalf("expected name to be %s, got %s", tc.output, name)
			}
		})
	}
}
