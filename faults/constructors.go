package faults

import (
	"fmt"

	flt "github.com/PlayerR9/go-error"
)

// FromString creates a new Fault with the given message and Unknown as its code.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func FromString(msg string) flt.Fault {
	return flt.New(Unknown, msg)
}

// FromStringf is like FromString but with a format string.
//
// Parameters:
//   - format: The format string of the fault.
//   - args: The arguments.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func FromStringf(format string, args ...any) flt.Fault {
	return flt.New(Unknown, fmt.Sprintf(format, args...))
}
