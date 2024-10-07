package faults

import (
	flt "github.com/PlayerR9/go-fault"
)

// LinesOf returns the fault's message and its embedding tower as a list of strings.
//
// Parameters:
//   - fault: The fault whose additional information are to be written.
//
// Returns:
//   - []string: The fault's additional information.
//
// An empty line is added between the message and the embedding tower. Also, a "dot" is
// added at the end of the message.
func LinesOf(fault flt.Fault) []string {
	if fault == nil {
		return nil
	}

	var lines []string

	lines = append(lines, ErrorOf(fault)+".")
	lines = append(lines, "")

	tmp := flt.InfoLines(fault)
	lines = append(lines, tmp...)

	return lines
}
