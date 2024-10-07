package faults

import (
	"fmt"

	flt "github.com/PlayerR9/go-fault"
)

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
