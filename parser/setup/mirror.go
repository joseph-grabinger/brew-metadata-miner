package setup

import "regexp"

// isDefaultMirrorPattern returns true if the given line
// matches the mirror pattern. It also returns the matches.
func isDefaultMirrorPattern(line string) (bool, []string) {
	regex := regexp.MustCompile(mirrorPattern)
	matches := regex.FindStringSubmatch(line)
	return len(matches) >= 2, matches
}
