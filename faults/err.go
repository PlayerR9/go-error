package faults

import (
	"fmt"

	flt "github.com/PlayerR9/go-fault"
)

// ErrFault is a fault that wraps a Go error.
type ErrFault struct {
	flt.Fault

	// Err is the error that occurred.
	Err error
}

// Embeds implements the flt.Fault interface.
func (e ErrFault) Embeds() flt.Fault {
	return e.Fault
}

// InfoLines implements the flt.Fault interface.
//
// Format:
//
//	"- Error: <error>"
//
// Where, <error> is the error that occurred. If no error was provided, "no error
// provided" is used instead.
func (e ErrFault) InfoLines() []string {
	lines := make([]string, 0, 1)

	if e.Err != nil {
		lines = append(lines, "- Error: "+e.Err.Error())
	} else {
		lines = append(lines, "- Error: no error provided")
	}

	return lines
}

// FromErr creates a new ErrFault.
//
// Parameters:
//   - err: The error that occurred.
//
// Returns:
//   - *ErrFault: The new ErrFault. Never returns nil.
func FromErr(err error) flt.Fault {
	base := flt.New(flt.UnknownCode, "something went wrong")

	return &ErrFault{
		Fault: base,
		Err:   err,
	}
}

// ErrPanic is an error that indicates that a panic occurred.
type ErrPanic struct {
	flt.Fault

	// Value is the value that was passed to panic.
	Value any
}

// Embeds implements the flt.Fault interface.
func (e ErrPanic) Embeds() flt.Fault {
	return e.Fault
}

// InfoLines implements the flt.Fault interface.
//
// Format:
//
//	"- Value: <value>"
func (e ErrPanic) InfoLines() []string {
	lines := make([]string, 0, 1)

	lines = append(lines, "- Value: "+fmt.Sprintf("%v", e.Value))

	return lines
}

// NewErrPanic creates a new ErrPanic.
//
// Parameters:
//   - value: The value that was passed to panic.
//
// Returns:
//   - *ErrPanic: A new ErrPanic. Never returns nil.
func NewErrPanic(value any) *ErrPanic {
	base := flt.WithLevel(flt.FATAL, flt.UnknownCode, "a panic occurred")

	return &ErrPanic{
		Fault: base,
		Value: value,
	}
}
