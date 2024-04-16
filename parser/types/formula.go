package types

import (
	"fmt"
	"log"
)

// Formula represents a formula from the brew package manager.
type Formula struct {
	// Name of the formula.
	Name string

	// Repository URL of the formula.
	RepoURL string

	// Archive URL of the formula.
	ArchiveURL string

	// License of the formula.
	License string

	// A list of the formula's dependencies.
	Dependencies []*Dependency
}

func (f *Formula) String() string {
	return fmt.Sprintf("%s\nRepo: %s\nArchive: %s\nLicense: %s\nDependencies: %v\n", f.Name, f.RepoURL, f.ArchiveURL, f.License, f.Dependencies)
}

// fromSourceFormula creates a formula from a source formula and evaluates the reopURL.
// It returns a pointer to the newly created formula.
func FromSourceFormula(sf *SourceFormula) *Formula {
	f := &Formula{
		Name:         sf.Name,
		License:      sf.formatLicense(),
		Dependencies: sf.Dependencies,
		ArchiveURL:   sf.URL,
	}

	repoURL, err := sf.extractRepoURL()
	if err != nil {
		log.Println(err)
		repoURL = ""
	}

	f.RepoURL = repoURL

	return f
}
