package setup

import (
	"regexp"

	"main/parser/types"
)

// cleanDependencySequence returns a cleaned slice of dependencies from a given sequence.
func cleanDependencySequence(sequence []string) []*types.Dependency {
	return cleanDepSequence(sequence, []skipSequence{}, 0)
}

// isBeginDependencySequence returns true if the given line
// is the beginning of a dependency sequence.
func isBeginDependencySequence(line string) bool {
	regex := regexp.MustCompile(beginDependencyPattern)
	return regex.MatchString(line)
}

// isEndDependencySequence returns true if the given line
// is the end of a dependency sequence.
func isEndDependencySequence(line string) bool {
	regex := regexp.MustCompile(endDependencyPatternNegated)
	return !regex.MatchString(line)
}
