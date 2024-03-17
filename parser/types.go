package parser

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// formula represents a formula from the brew package manager.
type formula struct {
	// name of the formula.
	name string

	// Repository URL of the formula.
	repoURL string

	// license of the formula.
	license string

	// A list of the formula's dependencies.
	dependencies []*dependency
}

func (f *formula) String() string {
	return fmt.Sprintf("%s\nRepo: %s\nLicense: %s\nDependencies: %v\n", f.name, f.repoURL, f.license, f.dependencies)
}

// fromSourceFormula creates a formula from a source formula and evaluates the reopURL.
// It returns a pointer to the newly created formula.
func fromSourceFormula(sf *sourceFormula) *formula {
	f := &formula{
		name:         sf.name,
		license:      sf.extractLicense(),
		dependencies: sf.dependencies,
	}

	repoURL, err := sf.extractRepoURL()
	if err != nil {
		log.Println(err)
		repoURL = ""
	}

	f.repoURL = repoURL

	return f
}

// sourceFormula represents a formula as found in the formula file.
type sourceFormula struct {
	// name of the formula.
	name string

	// homepage of the formula.
	homepage string

	// url of the formula.
	url string

	// mirror of the formula.
	mirror string

	// license of the formula.
	license string

	// List of the formula's dependencies.
	dependencies []*dependency

	// head of the formula.
	head *head
}

func (sf *sourceFormula) String() string {
	return fmt.Sprintf("%s\nHomepage: %s\nURL: %s\nMirror: %s\nLicense: %s\nDependencies: %v\nHead: %v\n", sf.name, sf.homepage, sf.url, sf.mirror, sf.license, sf.dependencies, sf.head)
}

// extractRepoURL returns the repository URL of the formula.
// It therfore inspects the URL, mirror and homepage fields of the formula.
func (sf *sourceFormula) extractRepoURL() (string, error) {
	var repoURL string

	// Use head if it exists.
	if sf.head != nil {
		return sf.head.url, nil
	}

	// Check homepage for known repository hosts.
	if m, repoURL := matchesKnownGitRepoHost(sf.homepage); m {
		return repoURL, nil
	}

	if strings.Contains(sf.homepage, "git.") { // TODO: Check if this is a good indicator and handle accordingly.
		log.Println("HOMEPAGE CONTAINS GIT: ", sf.homepage, sf.name)
	}

	if sf.url != "" {
		repoURL = sf.url
	} else if sf.mirror != "" {
		repoURL = sf.mirror
	} else {
		// Use homepage as fallback.
		repoURL = sf.homepage
	}

	if m, cleandedURL := matchesKnownGitRepoHost(repoURL); m {
		return cleandedURL, nil
	}

	if m, cleandedURL := matchesKnownGitArchiveHost(repoURL); m {
		return cleandedURL, nil
	}

	if strings.HasSuffix(repoURL, ".git") {
		return repoURL, nil
	}

	return "", fmt.Errorf("no repository URL found for formula: %s, repoURL: %s", sf.name, repoURL)
}

func (sf *sourceFormula) extractLicense() string {
	if sf.license == "" {
		return "pseudo"
	}

	license := strings.ReplaceAll(sf.license, "\"", "")

	re := regexp.MustCompile(`\b[\w.-]+\b`)
	matches := re.FindAllString(license, -1)
	if len(matches) == 1 {
		license = matches[0]
	} else if matches[0] == "any_of" || matches[0] == "one_of" {
		license = strings.Join(matches[1:], " or ")
	} else if matches[0] == "all_of" {
		license = strings.Join(matches[1:], " and ")
	}

	license = strings.ReplaceAll(license, "public_domain", "Public Domain")
	license = strings.ReplaceAll(license, "cannot_represent", "Cannot Represent")

	return license
}

// dependency represents a dependency of a formula.
type dependency struct {
	// name of the dependency.
	name string

	// depType is the type of the dependency.
	depType string
}

func (d *dependency) String() string {
	return fmt.Sprintf("{%s %s}", d.name, d.depType)
}

// head represents the head of a formula.
type head struct {
	// url of the head.
	url string

	// Version control system used.
	vcs string

	// dependencies of the head.
	dependencies []*dependency
}

func (h *head) String() string {
	return fmt.Sprintf("{%s %s, %s}", h.url, h.vcs, h.dependencies)
}
