package parser

import "fmt"

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

func (f *sourceFormula) String() string {
	return fmt.Sprintf("%s\nHomepage: %s\nURL: %s\nMirror: %s\nLicense: %s\nDependencies: %v\nHead: %v\n", f.name, f.homepage, f.url, f.mirror, f.license, f.dependencies, f.head)
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
}

func (h *head) String() string {
	return fmt.Sprintf("{%s %s}", h.url, h.vcs)
}
