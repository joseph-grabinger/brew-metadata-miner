package parser

import (
	"regexp"
	"strings"
)

// cleanLicenseSequence returns a cleaned string from a sequence.
func cleanLicenseSequence(sequence []string) string {
	for i, line := range sequence {
		// Remove comments.
		regex := regexp.MustCompile(`#.*$`)
		sequence[i] = regex.ReplaceAllString(line, "")
		// Remove whitespace.
		sequence[i] = strings.ReplaceAll(line, " ", "")
		// Remove tabs.
		sequence[i] = strings.ReplaceAll(line, "\t", "")
		// sequence[i] = strings.TrimSpace(line)
	} // TODO: Rework this.
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
