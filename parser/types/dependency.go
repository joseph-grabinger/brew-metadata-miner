package types

import (
	"fmt"
)

// Dependency represents a Dependency of a formula.
type Dependency struct {
	// Name of the dependency.
	Name string

	// DepType is the type of the dependency.
	DepType string
}

func (d *Dependency) String() string {
	return fmt.Sprintf("{%s %s}", d.Name, d.DepType)
}
