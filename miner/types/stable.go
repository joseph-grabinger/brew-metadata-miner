package types

import "fmt"

// Stable represents a stable version of a formula.
type Stable struct {
	// URL of the fomula's stable version.
	URL string

	// Dependencies of the stable version.
	Dependencies *Dependencies
}

func (s *Stable) String() string {
	return fmt.Sprintf("{%s, %s}", s.URL, s.Dependencies)
}
