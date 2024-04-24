package types

import (
	"fmt"
	"regexp"
	"strings"
)

// SourceFormula represents a formula as found in the formula file.
type SourceFormula struct {
	// Name of the formula.
	Name string

	// Homepage of the formula.
	Homepage string

	// Stable version of the formula.
	Stable *Stable

	// Mirror of the formula.
	Mirror string

	// License of the formula.
	License string

	// List of the formula's Dependencies.
	Dependencies *Dependencies

	// Head of the formula.
	Head *Head
}

func (sf *SourceFormula) String() string {
	return fmt.Sprintf("%s\nHomepage: %s\nStable: %s\nMirror: %s\nLicense: %s\nDependencies: %v\nHead: %v", sf.Name, sf.Homepage, sf.Stable, sf.Mirror, sf.License, sf.Dependencies, sf.Head)
}

// extractRepoURL returns the repository URL of the formula.
// It therfore inspects the URL, mirror and homepage fields of the formula.
func (sf *SourceFormula) extractRepoURL() string {
	var repoURL string

	// Use head if it exists.
	if sf.Head != nil {
		return sf.Head.URL
	}

	// Check homepage for known repository hosts.
	if m, repoURL := matchesKnownGitRepoHost(sf.Homepage); m {
		return repoURL
	}

	if sf.Stable != nil && sf.Stable.URL != "" {
		repoURL = sf.Stable.URL
	} else if sf.Mirror != "" {
		repoURL = sf.Mirror
	} else {
		return ""
	}

	if m, cleandedURL := matchesKnownGitRepoHost(repoURL); m {
		return cleandedURL
	}

	if m, cleandedURL := matchesKnownGitArchiveHost(repoURL); m {
		return cleandedURL
	}

	if strings.HasSuffix(repoURL, ".git") {
		return repoURL
	}

	return ""
}

func (sf *SourceFormula) formatLicense() string {
	if sf.License == "" {
		return "pseudo"
	}

	// Remove spaces and quotes.
	r := strings.NewReplacer(
		" ", "",
		"\"", "",
	)
	license := r.Replace(sf.License)

	// Remove unnecessary curly brackets.
	r = strings.NewReplacer(
		",{", ",",
		"]}", "]",
		",}", "}",
	)
	license = r.Replace(license)

	result := make([]rune, 0)
	sequence := make([]string, 0)
	word := make([]rune, 0)
	open, close := 0, 0
	operator := ""
	for _, r := range license {
		switch r {
		case ',':
			if len(word) > 0 {
				w := string(word)
				// Check for license exceptions.
				if operator != "" && strings.Contains(w, "=>{with:") {
					w = "(" + w + ")"
				}
				sequence = append(sequence, w)
				word = make([]rune, 0)
			}
		case '[':
			open++

			if len(sequence) > 0 {
				joined := []rune(strings.Join(sequence, operator))
				result = append(result, joined...)

				// Check if open bracket is needed.
				if open > 1 {
					result = append(result, []rune(operator+"(")...)
				}

				sequence = make([]string, 0)
			}

			if string(word) == "any_of:" || string(word) == "one_of:" {
				operator = " or "
			} else if string(word) == "all_of:" {
				operator = " and "
			}

			word = make([]rune, 0)
		case ']':
			close++

			if len(word) > 0 {
				w := string(word)
				// Check for license exceptions.
				if operator != "" && strings.Contains(w, "=>{with:") {
					w = "(" + w + ")"
				}
				sequence = append(sequence, w)
				word = make([]rune, 0)
			}

			joined := []rune(strings.Join(sequence, operator))
			result = append(result, joined...)
			sequence = make([]string, 0)

			// Check if close bracket is needed.
			if close > 1 {
				result = append(result, ')')
			}
		default:
			word = append(word, r)
		}
	}

	if len(word) > 0 {
		result = word
	}

	// Handle license exceptions.
	r = strings.NewReplacer(
		":public_domain", "Public Domain",
		":cannot_represent", "Cannot Represent",
		"=>", " ",
		":", " ",
		"{", "",
		"}", "",
	)
	return r.Replace(string(result))
}

// Known hosts for repository extraction.
const (
	// githubRepoPattern matches the URL of a Github repository.
	githubRepoPattern = `https://github.com/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)(/|\.git|\?.*)?$`

	// gitlabRepoPattern matches the URL of a Gitlab repository.
	gitlabRepoPattern = `https://gitlab.com/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)(/|\.git|\?.*)?$`

	// bitbucketRepoPattern matches the URL of a Bitbucket repository.
	bitbucketRepoPattern = `https://bitbucket.org/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)(/|\.git|\?.*)?$`

	// repoPattern represents a general repo pattern that matches the URL of any repository.
	repoPattern = `(https:\/\/[a-zA-Z0-9.-]+)\/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)`

	// githubArchivePattern matches the URL of a Github archive.
	githubArchivePattern = `(https://github.com/[a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+)\/(?:releases\/download|archive)\/.*`

	// gitlabArchivePattern matches the URL of a Gitlab archive.
	gitlabArchivePattern = `(https://gitlab.com/[a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+)\/(-\/archive|uploads)\/.*`

	//bitbucketArchivePattern matches the URL of a Bitbucket archive.
	bitbucketArchivePattern = `(https://bitbucket.org/[a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+)\/(downloads|get)\/.*`
)

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
