package faults

import (
	"fmt"

	flt "github.com/PlayerR9/go-fault"
)

// Join is a helper function that joins a list of faults into a single fault.
type JoinFault struct {
	flt.Fault

	// faults contains all the faults that have been joined.
	faults []flt.Fault
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (jf JoinFault) Embeds() flt.Fault {
	return jf.Fault
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Faults: <faults>"
func (jf JoinFault) InfoLines() []string {
	var lines []string

	for _, fault := range jf.faults {
		tmp := fault.InfoLines()
		lines = append(lines, tmp...)
	}

	return lines
}

// Join is a helper function that joins a list of faults into a single fault.
//
// Parameters:
//   - faults: The faults to join. May be nil.
//
// Returns:
//   - flt.Fault: The joined fault.
//
// This function returns nil if all the faults are nil.
func Join(faults ...flt.Fault) flt.Fault {
	// 1. Remove nil faults.
	var count int

	for _, fault := range faults {
		if fault != nil {
			count++
		}
	}

	if count == 0 {
		return nil
	}

	result := make([]flt.Fault, 0, count)

	for i := 0; i < len(faults); i++ {
		if faults[i] != nil {
			result = append(result, faults[i])
		}
	}

	// 2. Get the highest level of severity.
	highest := LevelOf(faults[0])

	for _, fault := range faults[1:] {
		level := LevelOf(fault)

		if level > highest {
			highest = level
		}
	}

	base := flt.WithLevel(highest, flt.FaultJoin, fmt.Sprintf("joined %d faults", len(faults)))

	js := &JoinFault{
		Fault:  base,
		faults: faults,
	}

	return js
}
