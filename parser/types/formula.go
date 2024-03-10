package types

import "fmt"

// TODO: Use the following representtion once the relevant meta data is stripped.
// Formula represents a formula from the brew package manager.
// type Formula struct {
// 	// Name of the formula.
// 	Name string

// 	// Repository URL of the formula.
// 	RepoURL string

// 	// License of the formula.
// 	License string

// 	// A list of the formula's dependencies.
// 	Dependencies []*Dependency
// }

// Formula represents a formula from the brew package manager.
type Formula struct {
	// Name of the formula.
	Name string

	// Homepage of the formula.
	Homepage string

	// URL of the formula.
	URL string

	// Mirror of the formula.
	Mirror string

	// License of the formula.
	License string

	// List of the formula's dependencies.
	Dependencies []*Dependency

	// Head of the formula.
	Head *Head
}

func (f *Formula) String() string {
	return fmt.Sprintf("%s\nHomepage: %s\nURL: %s\nMirror: %s\nLicense: %s\nDependencies: %v\nHead: %v\n", f.Name, f.Homepage, f.URL, f.Mirror, f.License, f.Dependencies, f.Head)
}

// Dependency represents a dependency of a formula.
type Dependency struct {
	// Name of the dependency.
	Name string

	// Type of the dependency.
	Type string
}

func (d *Dependency) String() string {
	return fmt.Sprintf("{%s %s}", d.Name, d.Type)
}

// Head represents the head of a formula.
type Head struct {
	// URL of the head.
	URL string

	// Version control system used.
	VCS string
}

func (h *Head) String() string {
	return fmt.Sprintf("{%s %s}", h.URL, h.VCS)
}
