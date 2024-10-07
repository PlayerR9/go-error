package fault

// StandardCode is a set of common, standard fault codes.
type StandardCode int

const (
	// Invalid specifies when faults are constructed in a way that is not expected.
	// This is used mostly for internal usage. However, users may use this code
	// when constructing faults in their own code.
	Invalid StandardCode = iota - 1

	// UnknownCode specifies when a fault is created but no code was specified.
	UnknownCode

	// FaultJoin is the fault code for when a fault is joined.
	FaultJoin

	// BadParameter specifies a broad category of faults that are caused by unexpected
	// parameters.
	BadParameter

	// OperationFailed specifies a broad category of faults that are returned by functions
	// when they fail.
	OperationFailed
)

var (
	// BadConstruction occurs when a fault does not embed *baseFault.
	BadConstruction FaultDescriber
)

func init() {
	BadConstruction = NewDescriptor(FATAL, Invalid, "fault does not implement *baseFault")
}
