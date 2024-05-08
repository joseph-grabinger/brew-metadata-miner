package types

import (
	"fmt"
	"strings"
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

// FormatPackageLine formats the formula as a package line.
// `0,"<package_manager>","<name>","<license>","<namespace>/<username>/<repository>","<stable_archive_url>","<system_requirement>"`
func (f *Formula) FormatPackageLine() string {
	return fmt.Sprintf("0\t\"brew\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\n", f.Name, f.License, f.RepoURL, f.ArchiveURL, f.SystemRequirement)
}

// FormatDependencyLine formats the formula as a dependency line.
// `1,"<package_manager>","<name>","<license>","<type>","<system_restriction>"`
func (f *Formula) FormatDependencyLine(dep *Dependency) string {
	depType := ""
	if len(dep.DepType) > 0 {
		depType = strings.Join(dep.DepType, ",")
	}
	return fmt.Sprintf("1\t\"brew\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\n", dep.Name, f.License, depType, dep.Restriction)
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
