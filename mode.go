package datacontext

import (
	"fmt"

	"github.com/mandelsoft/goutils/general"
)

// BuilderMode controls the handling of unset information in the
// builder configuration when calling the New method.
type BuilderMode int

const (
	// MODE_SHARED uses the default contexts for unset nested context types.
	MODE_SHARED BuilderMode = iota
	// MODE_DEFAULTED uses dedicated context instances configured with the
	// context type specific default registrations.
	MODE_DEFAULTED
	// MODE_EXTENDED uses dedicated context instances configured with
	// context type registrations extending the default registrations.
	MODE_EXTENDED
	// MODE_CONFIGURED uses dedicated context instances configured with the
	// context type registrations configured with the actual state of the
	// default registrations.
	MODE_CONFIGURED
	// MODE_INITIAL uses completely new contexts for unset nested context types
	// and initial registrations.
	MODE_INITIAL
)

const MULTI_REF = false

func (m BuilderMode) String() string {
	switch m {
	case MODE_SHARED:
		return "shared"
	case MODE_DEFAULTED:
		return "defaulted"
	case MODE_EXTENDED:
		return "extended"
	case MODE_CONFIGURED:
		return "configured"
	case MODE_INITIAL:
		return "initial"
	default:
		return fmt.Sprintf("(invalid %d)", m)
	}
}

func Mode(m ...BuilderMode) BuilderMode {
	return general.OptionalDefaulted(MODE_EXTENDED, m...)
}
