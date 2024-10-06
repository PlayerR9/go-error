package fault

import (
	"fmt"
	"time"
)

// Join is a helper function that joins a list of faults into a single fault.
type joinFault struct {
	// level indicates the severity level of the fault.
	level FaultLevel

	// msg informs about the nature of the fault.
	msg string

	// timestamp specifies the time when the fault occurred.
	timestamp time.Time

	// faults contains all the faults that have been joined.
	faults []Fault
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (jf joinFault) Embeds() Fault {
	return nil
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Faults: <faults>"
func (jf joinFault) InfoLines() []string {
	var lines []string

	for _, fault := range jf.faults {
		tmp := fault.InfoLines()
		lines = append(lines, tmp...)
	}

	return lines
}

// Error implements the Fault interface.
func (jf joinFault) Error() string {
	// [] () msg

	panic("not implemented")
}

// Level implements the Fault interface.
func (jf joinFault) Level() FaultLevel {
	return jf.level
}

// Timestamp implements the Fault interface.
func (jf joinFault) Timestamp() time.Time {
	return jf.timestamp
}

// Join is a helper function that joins a list of faults into a single fault.
//
// Parameters:
//   - faults: The faults to join. May be nil.
//
// Returns:
//   - Fault: The joined fault.
//
// This function returns nil if all the faults are nil.
func Join(faults ...Fault) Fault {
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

	result := make([]Fault, 0, count)

	for i := 0; i < len(faults); i++ {
		if faults[i] != nil {
			result = append(result, faults[i])
		}
	}

	// 2. Get the highest level of severity.
	highest := faults[0].Level()

	for _, fault := range faults[1:] {
		level := fault.Level()

		if level > highest {
			highest = level
		}
	}

	// 3. Get the oldest fault.
	min := faults[0].Timestamp()

	for i := 1; i < len(faults); i++ {
		timestamp := faults[i].Timestamp()

		if timestamp.Before(min) {
			min = timestamp
		}
	}

	js := &joinFault{
		faults:    faults,
		level:     highest,
		msg:       fmt.Sprintf("joined %d faults", len(faults)),
		timestamp: min,
	}

	return js
}
