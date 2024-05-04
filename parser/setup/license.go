package setup

import (
	"regexp"
	"strings"
)

// cleanLicenseSequence returns a cleaned string from a sequence.
func cleanLicenseSequence(sequence []string) string {
	// Remove leading license keyword.
	regex := regexp.MustCompile(licenseKeywordPattern)
	sequence[0] = regex.ReplaceAllString(sequence[0], "")
	for i := range sequence {
		// Remove comments.
		regex := regexp.MustCompile(commentPattern)
		sequence[i] = regex.ReplaceAllString(sequence[i], "")
		// Remove whitespace, tabs, and newlines.
		sequence[i] = strings.TrimSpace(sequence[i])
	}
	return strings.Join(sequence, "")
}

// isDefaultLicensePattern returns true if the given line
// matches the license pattern. It also returns the matches.
func isDefaultLicensePattern(line string) (bool, []string) {
	regex := regexp.MustCompile(licensePattern)
	matches := regex.FindStringSubmatch(line)
	return len(matches) >= 2, matches
}

// isBeginLicenseSequence returns true if the given line
// is the beginning of a license sequence.
func isBeginLicenseSequence(line string) bool {
	match, _ := regexp.MatchString(licenseKeywordPattern, line)
	return match && hasUnclosedBrackets(line)
}

// isEndLicenseSequence returns true if the given line
// is the end of a license sequence.
func isEndLicenseSequence(line string) bool {
	match, _ := regexp.MatchString(trailingCommaPattern, line)
	return hasUnopenedBrackets(line) && !match
}

// hasUnclosedBrackets returns true if the given line
// has more opening than closing brackets.
func hasUnclosedBrackets(line string) bool {
	open, close := countBrackets(line)
	return open > close
}

// hasUnopenedBrackets returns true if the given line
// has more closing than opening brackets.
func hasUnopenedBrackets(line string) bool {
	open, close := countBrackets(line)
	return open < close
}

// countBrackets returns the number of opening and
// closing square and curly brackets in the given string.
func countBrackets(s string) (open int, close int) {
	openCount, closeCount := 0, 0
	for _, char := range s {
		switch char {
		case '[', '{':
			openCount++
		case ']', '}':
			closeCount++
		}
	}
	return openCount, closeCount
}
