package types

import (
	"fmt"
	"log"
)

// Formula represents a Formula from the brew package manager.
type Formula struct {
	// Name of the formula.
	Name string

	// Repository URL of the formula.
	RepoURL string

	// License of the formula.
	License string

	// A list of the formula's Dependencies.
	Dependencies []*Dependency
}

func (f *Formula) String() string {
	return fmt.Sprintf("%s\nRepo: %s\nLicense: %s\nDependencies: %v\n", f.Name, f.RepoURL, f.License, f.Dependencies)
}

// fromSourceFormula creates a formula from a source formula and evaluates the reopURL.
// It returns a pointer to the newly created formula.
func FromSourceFormula(sf *SourceFormula) *Formula {
	f := &Formula{
		Name:         sf.Name,
		License:      sf.formatLicense(),
		Dependencies: sf.Dependencies,
	}

	repoURL, err := sf.extractRepoURL()
	if err != nil {
		log.Println(err)
		repoURL = ""
	}

	f.RepoURL = repoURL

	return f
}
