package setup

import (
	"fmt"
	"log"
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

		// Check for the tag.
		regex = regexp.MustCompile(tagPattern)
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
	log.Println("Init URL: ", url)
	url = strings.TrimSuffix(url, ".git")

	p := path.Join("tree", tag)
	return fmt.Sprintf("%s/%s", strings.TrimRight(url, "/"), strings.TrimLeft(p, "/"))
}

// isBeginURLSequence returns true if the given line
// is the beginning of a URL sequence.
func isBeginURLSequence(line string) bool {
	match, _ := regexp.MatchString(urlBeginPattern, line)
	return match
}

// isEndURLSequence returns true if the given line
// is the end of a URL sequence.
func isEndURLSequence(line string) bool {
	match, _ := regexp.MatchString(endPattern(2), line)
	return match
}
