package fault

import "fmt"

// New creates a new Fault given the code and its message. The timestamp is set to the
// current time.
//
// Parameters:
//   - code: The code of the fault.
//   - msg: The message of the fault.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
//
// The level of the fault is set to ERROR.
func New[C FaultCode](code C, msg string) Fault {
	descriptor := NewDescriptor(ERROR, code, msg)

	return descriptor.Init()
}

// WithLevel is like New but allows to specify the level of the fault.
//
// Parameters:
//   - level: The level of the fault.
//   - code: The code of the fault.
//   - msg: The message of the fault.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func WithLevel[C FaultCode](level FaultLevel, code C, msg string) Fault {
	descriptor := NewDescriptor(level, code, msg)

	return descriptor.Init()
}

// FromString creates a new Fault with the given message and Unknown as its code.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func FromString(msg string) Fault {
	return New(UnknownCode, msg)
}

// FromStringf is like FromString but with a format string.
//
// Parameters:
//   - format: The format string of the fault.
//   - args: The arguments.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func FromStringf(format string, args ...any) Fault {
	return New(UnknownCode, fmt.Sprintf(format, args...))
}
