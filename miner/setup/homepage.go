package setup

import "regexp"

// isDefaultHomepagePattern returns true if the given line
// matches the homepage pattern. It also returns the matches.
func isDefaultHomepagePattern(line string) (bool, []string) {
	regex := regexp.MustCompile(homepagePattern)
	matches := regex.FindStringSubmatch(line)
	return len(matches) >= 2, matches
}
