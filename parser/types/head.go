package types

import (
	"fmt"
)

// Head represents the Head of a formula.
type Head struct {
	// URL of the head.
	URL string

	// Version control system used.
	VCS string

	// Dependencies of the head.
	Dependencies []*Dependency
}

func (h *Head) String() string {
	return fmt.Sprintf("{%s %s, %s}", h.URL, h.VCS, h.Dependencies)
}
