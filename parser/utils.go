package parser

import (
	"regexp"
	"strings"
)

// cleanDependencySequence returns a cleaned [][]string from a sequence.
func cleanDependencySequence(sequence []string) interface{} {
	res := make([][]string, 0)
	for i := range sequence {
		// Check for system dependency.
		regex := regexp.MustCompile(`uses_from_macos\s+"([^"]+)"`)
		nameMatches := regex.FindStringSubmatch(sequence[i])
		if len(nameMatches) >= 2 {
			res = append(res, []string{nameMatches[1], "system"})
			continue
		}

		// Check for the dependency name.
		regex = regexp.MustCompile(`depends_on\s+"([^"]+)"`)
		nameMatches = regex.FindStringSubmatch(sequence[i])
		if len(nameMatches) < 2 {
			continue
		}

		// Check for the dependency type.
		regex = regexp.MustCompile(`=>\s*:(\w+)`)
		typeMatches := regex.FindStringSubmatch(sequence[i])
		if len(typeMatches) >= 2 {
			res = append(res, []string{nameMatches[1], typeMatches[1]})
		} else {
			res = append(res, []string{nameMatches[1]})
		}
	}
	return res
}

// isBeginDependencySequence returns true if the given line
// is the beginning of a dependency sequence.
func isBeginDependencySequence(line string) bool {
	regex := regexp.MustCompile(`^(\s{2}|\t)(depends_on|uses_from_macos)\s+"[^"]+"`)
	return regex.MatchString(line)
}

// isEndDependencySequence returns true if the given line
// is the end of a dependency sequence.
func isEndDependencySequence(line string) bool {
	regex := regexp.MustCompile(`^(\s{2}|\t)(depends_on|uses_from_macos)\s+|^[\s\t]*$`)
	return !regex.MatchString(line)
}

// cleanLicenseSequence returns a cleaned string from a sequence.
func cleanLicenseSequence(sequence []string) interface{} {
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
