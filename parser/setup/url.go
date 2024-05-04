package setup

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"main/parser/types"
)

// cleanURLSequence returns a cleaned string from a sequence.
func cleanURLSequence(sequence []string) *types.Stable {
	if len(sequence) == 1 {
		return &types.Stable{URL: sequence[0]}
	}

	stable := &types.Stable{}
	var index int
	for i := range sequence {
		// Check for the URL.
		if stable.URL == "" {
			regex := regexp.MustCompile(stableUrlPattern)
			urlMatches := regex.FindStringSubmatch(sequence[i])
			if len(urlMatches) >= 3 {
				stable.URL = urlMatches[1]
				if urlMatches[2] == "," {
					continue
				} else {
					index = i
					break
				}
			}
		}

		// Check for the tag.
		regex := regexp.MustCompile(tagPattern)
		tagMatches := regex.FindStringSubmatch(sequence[i])
		if len(tagMatches) >= 2 {
			stable.URL = formatURL(stable.URL, tagMatches[1])
			index = i
			break
		}
	}

	// Initialize skips for resources and patches.
	skips := skips{
		{begin: stableResourcePattern, end: endPattern(4)},
		{begin: stablePatchPattern, end: endPattern(4)},
	}

	// Check for dependencies.
	deps := cleanDepSequence(sequence[index+1:], skips, 1)
	stable.Dependencies = deps
	return stable
}

// formatURL joines the given url with the "tree" literal and given tag.
func formatURL(url, tag string) string {
	url = strings.TrimSuffix(url, ".git")

	p := path.Join("tree", tag)
	return fmt.Sprintf("%s/%s", strings.TrimRight(url, "/"), strings.TrimLeft(p, "/"))
}

// isDefaultURLPattern returns true if the given line
// is the default pattern for the URL field. In other words, the URL field
// can be extracted from a single line.
// The function returns a boolean indicating if the line matches the default pattern,
// and a slice containing the matches if the line matches the default pattern.
func isDefaultURLPattern(line string) (bool, []string) {
	regex := regexp.MustCompile(urlPattern)
	matches := regex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return false, nil
	}

	rem := strings.ReplaceAll(line, matches[0], "")
	// No trailing comma indicates default pattern.
	if !(len(rem) > 0 && rem[0] == ',') {
		return true, matches
	}

	// Check for tag within the same line.
	regex = regexp.MustCompile(tagExtractPattern)
	if tagMatches := regex.FindStringSubmatch(rem); len(tagMatches) >= 2 {
		return true, []string{
			fmt.Sprintf("%s, %s", matches[0], tagMatches[0]),
			formatURL(matches[1], tagMatches[1]),
		}
	}

	// The strings "using:" and "revision:" indicate default pattern.
	if strings.Contains(rem, "using:") || strings.Contains(rem, "revision:") {
		return true, matches
	}
	return false, nil
}

// isBeginURLSequence returns true if the given line
// is the beginning of a URL sequence.
func isBeginURLSequence(line string) bool {
	match, _ := regexp.MatchString(urlBeginPattern, line)
	return match && !(strings.Contains(line, "tag:") || strings.Contains(line, "using:") || strings.Contains(line, "revision:"))
}

// isEndURLSequence returns true if the given line
// is the end of a URL sequence.
func isEndURLSequence(line string) bool {
	endMatch, _ := regexp.MatchString(endPattern(2), line)
	revMatch, _ := regexp.MatchString(revisionPattern, line)
	return endMatch || revMatch
}
