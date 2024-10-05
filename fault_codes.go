package errs

// FaultCode is the interface that all fault codes should implement.
type FaultCode interface {
	~int

	// String returns the string representation of the fault code.
	//
	// Returns:
	//   - string: The string representation of the fault code.
	String() string
}

// StdFaultCode is the standard set of fault codes.
type StdFaultCode int

const (
	// OperationFail is the fault code for when an operation fails
	// for some reason.
	OperationFail StdFaultCode = iota
)
