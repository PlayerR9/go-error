package faults

// StdFaultCode is the standard set of fault codes.
type StdFaultCode int

const (
	// Unknown is the fault code for when an unknown error occurs.
	Unknown StdFaultCode = iota - 1

	// FaultJoin is the fault code for when a fault is joined.
	FaultJoin

	// OperationFail is the fault code for when an operation fails
	// for some reason.
	OperationFail
)
