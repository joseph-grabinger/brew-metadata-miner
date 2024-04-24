package types

import (
	"fmt"
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

	// System requirement of the formula.
	SystemRequirement string
}

func (f *Formula) String() string {
	return fmt.Sprintf("%s\nRepo: %s\nArchive: %s\nLicense: %s\nDependencies: %v\nSystemRequirements: %s\n", f.Name, f.RepoURL, f.ArchiveURL, f.License, f.Dependencies, f.SystemRequirement)
}

// FormatRepoLine formats the formula as a repository line.
// `0,"<namespace>/<username>/<repository>","<license>","<stablearchiveurl>"`
func (f *Formula) FormatRepoLine() string {
	return fmt.Sprintf("0\t\"%s\"\t\"%s\"\t\"%s\"\n", f.RepoURL, f.License, f.ArchiveURL)
}

// FormatDependencyLine formats the formula as a dependency line.
// `1,"<license>","<namespace>/<username>/<repository>","<stablearchiveurl>","<type>","<packagemanager>","<name>","<systemrestriction>"`
func (f *Formula) FormatDependencyLine(dep *Dependency) string {
	return fmt.Sprintf("1\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\n", f.License, f.RepoURL, f.ArchiveURL, dep.DepType, "brew", dep.Name, dep.Restriction)
}

// fromSourceFormula creates a formula from a source formula and evaluates the reopURL.
// It returns a pointer to the newly created formula.
func FromSourceFormula(sf *SourceFormula) *Formula {
	f := &Formula{
		Name:       sf.Name,
		License:    sf.formatLicense(),
		ArchiveURL: sf.Stable.URL,
	}

	if sf.Dependencies != nil {
		f.Dependencies = sf.Dependencies.Lst
		f.SystemRequirement = sf.Dependencies.SystemRequirements
	} else {
		f.Dependencies = []*Dependency{}
	}

	if sf.Stable.Dependencies != nil {
		f.Dependencies = append(f.Dependencies, sf.Stable.Dependencies.Lst...)

		if sf.Stable.Dependencies.SystemRequirements != "" {
			if f.SystemRequirement != "" {
				// Join the system requirements.
				f.SystemRequirement = fmt.Sprintf("%s, %s", f.SystemRequirement, sf.Stable.Dependencies.SystemRequirements)
			} else {
				f.SystemRequirement = sf.Stable.Dependencies.SystemRequirements
			}
		}
	}

	f.RepoURL = sf.extractRepoURL()

	return f
}
