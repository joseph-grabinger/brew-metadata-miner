package parser

import (
	"regexp"
	"strings"
)

// cleanLicenseSequence returns a cleaned string from a sequence.
func cleanLicenseSequence(sequence []string) string {
	// Remove leading license keyword.
	regex := regexp.MustCompile(`\s*license\s*`)
	sequence[0] = regex.ReplaceAllString(sequence[0], "")
	for i := range sequence {
		// Remove comments.
		regex := regexp.MustCompile(`#.*$`)
		sequence[i] = regex.ReplaceAllString(sequence[i], "")
		// Remove whitespace, tabs, and newlines.
		sequence[i] = strings.TrimSpace(sequence[i])
	}
	return strings.Join(sequence, "")
}

// isBeginLicenseSequence returns true if the given line
// is the beginning of a license sequence.
func isBeginLicenseSequence(line string) bool {
	match, _ := regexp.MatchString(`^\s*license`, line)
	return match && hasUnclosedBrackets(line)
}

// hasUnclosedBrackets returns true if the given line
// has more opening than closing square brackets.
func hasUnclosedBrackets(line string) bool {
	open, close := countBrackets(line)
	return open > close
}

// hasUnopenedBrackets returns true if the given line
// has more closing than opening square brackets.
func hasUnopenedBrackets(line string) bool {
	open, close := countBrackets(line)
	return open < close
}

// countBrackets returns the number of opening and
// closing square brackets in the given string.
func countBrackets(s string) (open int, close int) {
	openCount, closeCount := 0, 0
	for _, char := range s {
		switch char {
		case '[':
			openCount++
		case ']':
			closeCount++
		}
	}
	return openCount, closeCount
}
