package setup

import (
	"regexp"

	"main/miner/types"
)

// cleanHeadSequence returns a cleaned head from a sequence.
func cleanHeadSequence(sequence []string) *types.Head {
	if len(sequence) == 1 {
		return &types.Head{URL: sequence[0]}
	}

	head := &types.Head{}
	var index int
	for i := range sequence {
		// Check for the URL.
		regex := regexp.MustCompile(headBlockURLPattern)
		matches := regex.FindStringSubmatch(sequence[i])
		if len(matches) >= 2 {
			head.URL = matches[1]
			index = i
			break
		}
	}

	// Initialize skips for resources and patches.
	skips := skips{
		{begin: blockResourcePattern, end: endPattern(4)},
		{begin: blockPatchPattern, end: endPattern(4)},
	}

	deps := cleanDepSequence(sequence[index+1:], skips, 1)
	head.Dependencies = deps.Lst
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
