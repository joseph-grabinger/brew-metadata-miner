package types

import (
	"fmt"
)

// Dependency represents a dependency of a formula.
type Dependency struct {
	// Name of the dependency.
	Name string

	// DepType is the type of the dependency.
	DepType string

	// (System) restirction for the dependency.
	Restriction string
}

func (d *Dependency) String() string {
	return fmt.Sprintf("{%s %s %s}", d.Name, d.DepType, d.Restriction)
}

func (d *Dependency) Id() string {
	return fmt.Sprintf("%s,%s", d.Name, d.DepType)
}
