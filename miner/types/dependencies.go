package types

import "fmt"

// Dependecies represents a formula's dependencies.
// This struct is used to frther store the formula's system requirements,
// since these are defined in the same section of a formula's definiton.
type Dependencies struct {
	// List of dependencies.
	Lst []*Dependency

	// Formula's system requirements.
	SystemRequirements string
}

func (d *Dependencies) String() string {
	return fmt.Sprintf("{%v\n, SystemRequirements: %s}", d.Lst, d.SystemRequirements)
}
