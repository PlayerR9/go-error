package faults

import (
	"fmt"

	flt "github.com/PlayerR9/go-error/fault"
)

// FaultErr is a fault that wraps a Go error.
type FaultErr struct {
	flt.Fault

	// Err is the error that occurred.
	Err error
}

// Embeds implements the Fault interface.
func (e FaultErr) Embeds() flt.Fault {
	return e.Fault
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Error: <error>"
//
// Where, <error> is the error that occurred. If no error was provided, "no error
// provided" is used instead.
func (e FaultErr) InfoLines() []string {
	lines := make([]string, 0, 1)

	if e.Err != nil {
		lines = append(lines, "- Error: "+e.Err.Error())
	} else {
		lines = append(lines, "- Error: no error provided")
	}

	return lines
}

// NewFaultErr creates a new FaultErr.
//
// Parameters:
//   - err: The error that occurred.
//
// Returns:
//   - *FaultErr: The new FaultErr. Never returns nil.
func NewFaultErr(err error) flt.Fault {
	base := flt.New(Unknown, "something went wrong")

	return &FaultErr{
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

// Embeds implements the Fault interface.
func (e ErrPanic) Embeds() flt.Fault {
	return e.Fault
}

// InfoLines implements the Fault interface.
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
	base := flt.WithLevel(flt.FATAL, Unknown, "a panic occurred")

	return &ErrPanic{
		Fault: base,
		Value: value,
	}
}
