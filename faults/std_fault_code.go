package faults

// StdFaultCode is the standard set of fault codes.
type StdFaultCode int

const (
	// Unknown is the fault code for when an unknown error occurs.
	Unknown StdFaultCode = iota

	// OperationFail is the fault code for when an operation fails
	// for some reason.
	OperationFail
)
