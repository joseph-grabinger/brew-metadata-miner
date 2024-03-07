package parser

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

	// Archive of the formula.
	Archive *Archive

	// Head of the formula.
	Head *Head
}

// Dependency represents a dependency of a formula.
type Dependency struct {
	// Name of the dependency.
	Name string

	// Type of the dependency.
	Type string
}

// Archive represents a stable archive of a formula.
type Archive struct {
	// URL of the archive.
	URL string

	// Tag of the archive.
	Tag string
}

// Head represents the head of a formula.
type Head struct {
	// URL of the head.
	URL string

	// Version control system used.
	VCS string
}
