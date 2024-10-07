package fault

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Faulter interface {
}

// BaseFault is the base implementation of the Fault interface.
type BaseFault struct {
	// descriptor is the root information of any fault.
	descriptor FaultDescriber

	// timestamp specifies the time when the fault occurred.
	timestamp time.Time

	// suggestions describes one or more possible solutions or actions that can be taken
	// to resolve the fault.
	suggestions []string

	// stack_trace is the stack trace of the fault.
	stack_trace []string

	// context is the context of the fault.
	context map[string]any
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (bf BaseFault) Embeds() Fault {
	return nil
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"Occurred at: <timestamp>"
//
// "Suggestions:"
// "- <suggestion>"
// "- ..."
// "Stack trace:"
// "- <stack trace>"
//
// Where:
//   - <timestamp>: The time when the fault occurred.
//   - <suggestion>: One or more possible solutions or actions that can be taken
//     to resolve the fault.
//   - <stack trace>: The stack trace of the fault.
func (bf BaseFault) InfoLines() []string {
	var lines []string

	if !bf.timestamp.IsZero() {
		lines = append(lines, "Occurred at: "+bf.timestamp.String())
	}

	if len(bf.suggestions) > 0 {
		lines = append(lines, "Suggestions:")

		for _, suggestion := range bf.suggestions {
			lines = append(lines, "- "+suggestion)
		}
	}

	if len(bf.context) > 0 {
		lines = append(lines, "Context:")

		for k, v := range bf.context {
			lines = append(lines, fmt.Sprintf("- %s: %v", k, v))
		}
	}

	if len(bf.stack_trace) > 0 {
		lines = append(lines, "Stack trace:")

		trace := make([]string, len(bf.stack_trace), len(bf.stack_trace)+1)
		copy(trace, bf.stack_trace)
		trace = append(trace, "")

		slices.Reverse(trace)

		lines = append(lines, "- "+strings.Join(trace, " <- "))
	}

	return lines
}

func (bf *BaseFault) AppendFrame(frame string) bool {
	if bf == nil {
		return false
	}

	bf.stack_trace = append(bf.stack_trace, frame)

	return true
}

func (bf BaseFault) Error() string {
	return bf.descriptor.String()
}
