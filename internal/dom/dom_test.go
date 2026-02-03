package dom_test

import (
	"strings"
	"testing"

	"github.com/dhaef/job-scraper/internal/dom"
	"golang.org/x/net/html"
)

func TestFindElementsByClasses(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		classes   []string
		outputLen int
	}{
		{
			name:      "should find all elements",
			input:     "<div><p class=\"foo\">f</p><p class=\"bar\">b</p><p class=\"baz\">c</p><p class=\"foo\">f</p></div>",
			outputLen: 3,
			classes:   []string{"foo", "bar"},
		},
		{
			name:      "should find no elements",
			input:     "<div><p class=\"foo\">f</p class=\"bar\"><p>b</p><p class=\"baz\">c</p></div>",
			outputLen: 0,
			classes:   []string{"food"},
		},
		{
			name:      "should find single element",
			input:     "<div><p class=\"foo\">f</p class=\"bar\"><p>b</p><p class=\"baz\">c</p></div>",
			outputLen: 1,
			classes:   []string{"foo"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("failed to parse input, %v", err)
			}

			nodes := dom.FindElementsByClasses(doc, tc.classes)

			if len(nodes) != tc.outputLen {
				t.Fatalf("expected nodes length to be %d, got %d", tc.outputLen, len(nodes))
			}
		})
	}
}

func TestFindElementsByTagName(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		tagName   string
		outputLen int
	}{
		{
			name:      "should find all elements",
			input:     "<div><p class=\"foo\">f</p><p class=\"bar\">b</p></div>",
			outputLen: 2,
			tagName:   "p",
		},
		{
			name:      "should find no elements",
			input:     "<div><p class=\"foo\">f</p class=\"bar\"><p>b</p><p class=\"baz\">c</p></div>",
			outputLen: 0,
			tagName:   "a",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("failed to parse input, %v", err)
			}

			nodes := dom.FindElementsByTagName(doc, tc.tagName)

			if len(nodes) != tc.outputLen {
				t.Fatalf("expected nodes length to be %d, got %d", tc.outputLen, len(nodes))
			}
		})
	}
}

func TestFindLabelByClass(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		output     string
		labels     map[string][]string
		expectedOk bool
	}{
		{
			name:       "should find company label",
			input:      "<div><p class=\"foo\">f</p><p class=\"bar\">b</p></div>",
			output:     "company",
			expectedOk: true,
			labels: map[string][]string{
				"company": {"foo"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("failed to parse input, %v", err)
			}

			label, ok := dom.FindLabelByClass(doc, tc.labels)

			if ok != tc.expectedOk {
				t.Fatalf("expected ok to be %t, got %t", tc.expectedOk, ok)
			}

			if label != tc.output {
				t.Fatalf("expected label to be %s, got %s", tc.output, label)
			}
		})
	}
}

func TestGetTextContent(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "should return all nested text",
			input:  "<div><p class=\"foo\"><strong>hello:</strong> world</p></div>",
			output: "hello: world",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("failed to parse input, %v", err)
			}

			text := dom.GetTextContent(doc)

			if text != tc.output {
				t.Fatalf("expected text to be %s, got %s", tc.output, text)
			}
		})
	}
}

func TestGetHref(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
		exists bool
	}{
		{
			name:   "should return valid href",
			input:  `<a href="https://test.com">hello</a>`,
			output: "https://test.com",
			exists: true,
		},
		{
			name:   "should return false",
			input:  "<p>hello</p>",
			output: "",
			exists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("failed to parse input, %v", err)
			}

			href, ok := dom.GetHref(doc)
			if ok != tc.exists {
				t.Fatalf("expected exists to equal %t, got %t", tc.exists, ok)
			}

			if href != tc.output {
				t.Fatalf("expected href to be %s, got %s", tc.output, href)
			}
		})
	}
}
