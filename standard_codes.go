package fault

import "fmt"

// StandardCode is a set of common, standard fault codes.
type StandardCode int

const (
	// Invalid specifies when faults are constructed in a way that is not expected.
	// This is used mostly for internal usage. However, users may use this code
	// when constructing faults in their own code.
	Invalid StandardCode = iota - 1

	// Unknown specifies when a fault is created but no code was specified.
	Unknown

	// BadParameter specifies a broad category of faults that are caused by unexpected
	// parameters.
	BadParameter
)

var (
	// BadConstruction occurs when a fault does not embed *baseFault.
	BadConstruction FaultDescriber
)

func init() {
	BadConstruction = NewDescriptor(FATAL, Invalid, "fault does not implement *baseFault")
}

// NewBadParameter creates a new BadParameter fault.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func NewBadParameter(msg string) Fault {
	desc := NewDescriptor(ERROR, BadParameter, msg)

	return desc.Init()
}

// NewNilParameter creates a new BadParameter fault.
//
// Parameters:
//   - param_name: The name of the parameter.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func NewNilParameter(param_name string) Fault {
	desc := NewDescriptor(ERROR, BadParameter, fmt.Sprintf("parameter (%q) must be non-nil", param_name))

	return desc.Init()
}
