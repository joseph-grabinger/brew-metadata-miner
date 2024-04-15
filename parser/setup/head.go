package setup

import (
	"regexp"

	"main/parser/types"
)

// cleanHeadSequence returns a cleaned []string from a sequence.
func cleanHeadSequence(sequence []string) interface{} {
	head := &types.Head{Dependencies: make([]*types.Dependency, 0)}
	for i := range sequence {
		// Check for the URL.
		regex := regexp.MustCompile(headBlockURLPattern)
		matches := regex.FindStringSubmatch(sequence[i])
		if len(matches) >= 2 {
			head.URL = matches[1]

			// Check for the VCS.
			regex = regexp.MustCompile(headVCSPattern)
			matches = regex.FindStringSubmatch(sequence[i])
			if len(matches) >= 2 {
				head.VCS = matches[1]
			}
		}

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
			dep.DepType = typeMatches[1]
		}
		head.Dependencies = append(head.Dependencies, dep)
	}
	return head
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
	regex := regexp.MustCompile(endHeadPattern)
	return regex.MatchString(line)
}