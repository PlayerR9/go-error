package errs

import "time"

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
	// Unknown is the fault code for when an unknown error occurs.
	Unknown StdFaultCode = iota

	// OperationFail is the fault code for when an operation fails
	// for some reason.
	OperationFail
)

// FromString creates a new FaultBase struct with the given message and
// Unknown as the code. The timestamp is set to the current time.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - *FaultBase: The new FaultBase struct. Never returns nil.
func FromString(msg string) *FaultBase[StdFaultCode] {
	return &FaultBase[StdFaultCode]{
		level:     ERROR,
		code:      Unknown,
		msg:       msg,
		timestamp: time.Now(),
	}
}

////////////////////////////////////////////////////////////////////

type ErrPanic struct {
	FaultBase[StdFaultCode]
	Value any
}

// func (e ErrPanic) WriteInfo(w)
