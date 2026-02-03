package job

import "strings"

var LabelMap = map[string][]string{
	"company":    {"company"},
	"location":   {"location"},
	"salary":     {"salary", "salary range", "pay rate"},
	"department": {"department", "business unit", "business area", "business line", "business class", "function", "group"},
	"posted":     {"posted", "date", "published"},
}

func NormalizeMetadataLabel(label string) (string, bool) {
	if containsAny(label, LabelMap["company"]) {
		return "company", true
	}

	if containsAny(label, LabelMap["location"]) {
		return "location", true
	}

	if containsAny(label, LabelMap["salary"]) {
		return "salary", true
	}

	if containsAny(label, LabelMap["department"]) {
		return "department", true
	}

	if containsAny(label, LabelMap["posted"]) {
		return "posted", true
	}

	return "", false
}

func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
