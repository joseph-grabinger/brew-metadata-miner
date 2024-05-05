package setup

import (
	"regexp"
	"slices"

	"main/miner/types"
)

// cleanHeadSequence returns a cleaned head from a sequence.
func cleanHeadSequence(sequence []string) *types.Head {
	if len(sequence) == 1 {
		return &types.Head{URL: sequence[0]}
	}

	head := &types.Head{Dependencies: make([]*types.Dependency, 0)}
	for i := range sequence {
		// Check for the URL.
		regex := regexp.MustCompile(headBlockURLPattern)
		matches := regex.FindStringSubmatch(sequence[i])
		if len(matches) >= 2 {
			head.URL = matches[1]
		}

		// TODO use dep_utils.go to clean dependencies OR remove them entirely.
		// Check for dependencies.
		regex = regexp.MustCompile(dependencyKeywordPattern)
		matches = regex.FindStringSubmatch(sequence[i])
		if len(matches) < 2 {
			continue
		}

		dep := &types.Dependency{Name: matches[1]}

		// Check for the dependency type.
		regex = regexp.MustCompile(dependencyTypePattern)
		typeMatches := regex.FindStringSubmatch(sequence[i])
		if len(typeMatches) >= 2 {
			dep.DepType = slices.DeleteFunc(typeMatches[1:], func(s string) bool {
				return s == ""
			})
		}
		head.Dependencies = append(head.Dependencies, dep)
	}
	return head
}

// isDefaultHeadPattern returns true if the given line
// matches the head pattern. It also returns the matches.
func isDefaultHeadPattern(line string) (bool, []string) {
	regex := regexp.MustCompile(headURLPattern)
	matches := regex.FindStringSubmatch(line)
	return len(matches) >= 2, matches
}

// isBeginHeadSequence returns true if the given line
// is the beginning of a head sequence.
func isBeginHeadSequence(line string) bool {
	regex := regexp.MustCompile(beginHeadPattern)
	return regex.MatchString(line)
}

// isEndHeadSequence returns true if the given line
// is the end of a head sequence.
func isEndHeadSequence(line string) bool {
	regex := regexp.MustCompile(endPattern(2))
	return regex.MatchString(line)
}
