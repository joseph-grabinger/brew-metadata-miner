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
		regex := regexp.MustCompile(systemDependencyPattern)
		nameMatches := regex.FindStringSubmatch(sequence[i])
		if len(nameMatches) >= 2 {
			res = append(res, []string{nameMatches[1], "system"})
			continue
		}

		// Check for the dependency name.
		regex = regexp.MustCompile(dependencyKeywordPattern)
		nameMatches = regex.FindStringSubmatch(sequence[i])
		if len(nameMatches) < 2 {
			continue
		}

		// Check for the dependency type.
		regex = regexp.MustCompile(dependencyTypePattern)
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
	regex := regexp.MustCompile(beginDependencyPattern)
	return regex.MatchString(line)
}

// isEndDependencySequence returns true if the given line
// is the end of a dependency sequence.
func isEndDependencySequence(line string) bool {
	regex := regexp.MustCompile(endDependencyPattern)
	return !regex.MatchString(line)
}

// cleanHeadSequence returns a cleaned []string from a sequence.
func cleanHeadSequence(sequence []string) interface{} {
	head := &head{dependencies: make([]*dependency, 0)}
	for i := range sequence {
		// Check for the URL.
		regex := regexp.MustCompile(headBlockURLPattern)
		matches := regex.FindStringSubmatch(sequence[i])
		if len(matches) >= 2 {
			head.url = matches[1]

			// Check for the VCS.
			regex = regexp.MustCompile(headVCSPattern)
			matches = regex.FindStringSubmatch(sequence[i])
			if len(matches) >= 2 {
				head.vcs = matches[1]
			}
		}

		// Check for dependencies.
		regex = regexp.MustCompile(dependencyKeywordPattern)
		matches = regex.FindStringSubmatch(sequence[i])
		if len(matches) < 2 {
			continue
		}

		dep := &dependency{name: matches[1]}

		// Check for the dependency type.
		regex = regexp.MustCompile(dependencyTypePattern)
		typeMatches := regex.FindStringSubmatch(sequence[i])
		if len(typeMatches) >= 2 {
			dep.depType = typeMatches[1]
		}
		head.dependencies = append(head.dependencies, dep)
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

// cleanLicenseSequence returns a cleaned string from a sequence.
func cleanLicenseSequence(sequence []string) interface{} {
	// Remove leading license keyword.
	regex := regexp.MustCompile(licenseKeywordPattern)
	sequence[0] = regex.ReplaceAllString(sequence[0], "")
	for i := range sequence {
		// Remove comments.
		regex := regexp.MustCompile(commentPattern)
		sequence[i] = regex.ReplaceAllString(sequence[i], "")
		// Remove whitespace, tabs, and newlines.
		sequence[i] = strings.TrimSpace(sequence[i])
	}
	return strings.Join(sequence, "")
}

// isBeginLicenseSequence returns true if the given line
// is the beginning of a license sequence.
func isBeginLicenseSequence(line string) bool {
	match, _ := regexp.MatchString(licenseKeywordPattern, line)
	return match && hasUnclosedBrackets(line)
}

// hasUnclosedBrackets returns true if the given line
// has more opening than closing brackets.
func hasUnclosedBrackets(line string) bool {
	open, close := countBrackets(line)
	return open > close
}

// hasUnopenedBrackets returns true if the given line
// has more closing than opening brackets.
func hasUnopenedBrackets(line string) bool {
	open, close := countBrackets(line)
	return open < close
}

// countBrackets returns the number of opening and
// closing square and curly brackets in the given string.
func countBrackets(s string) (open int, close int) {
	openCount, closeCount := 0, 0
	for _, char := range s {
		switch char {
		case '[', '{':
			openCount++
		case ']', '}':
			closeCount++
		}
	}
	return openCount, closeCount
}

// matchesKnownGitRepoHost checks if the given url matches a known git repository host pattern.
// If true, it returns the matched repository url.
func matchesKnownGitRepoHost(url string) (bool, string) {
	githubRe := regexp.MustCompile(githubRepoPattern)
	gitlabRe := regexp.MustCompile(gitlabRepoPattern)
	bitBucketRe := regexp.MustCompile(bitbucketRepoPattern)

	if !(githubRe.MatchString(url) || gitlabRe.MatchString(url) || bitBucketRe.MatchString(url)) {
		return false, ""
	}

	matches := regexp.MustCompile(repoPattern).FindStringSubmatch(url)

	return true, matches[0] + ".git"
}

// matchesKnownGitArchiveHost checks if the given url matches a known git archive host pattern.
// If true, it returns the matched repository url.
func matchesKnownGitArchiveHost(url string) (bool, string) {
	githubRe := regexp.MustCompile(githubArchivePattern)
	if githubRe.MatchString(url) {
		matches := githubRe.FindStringSubmatch(url)
		return true, matches[1] + ".git"
	}

	gitlabRe := regexp.MustCompile(gitlabArchivePattern)
	if gitlabRe.MatchString(url) {
		matches := gitlabRe.FindStringSubmatch(url)
		return true, matches[1] + ".git"
	}

	bitbucketRe := regexp.MustCompile(bitbucketArchivePattern)
	if bitbucketRe.MatchString(url) {
		matches := bitbucketRe.FindStringSubmatch(url)
		return true, matches[1] + ".git"
	}

	return false, ""
}
