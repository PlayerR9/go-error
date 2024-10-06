package errs

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

func (jf joinFault) WriteInfo(w Writer) (int, Fault) {
	for _, fault := range jf.faults {
		WriteFault(nil, fault)
	}

	return 0, nil
}

func (jf joinFault) Error() string {
	// [] () msg

	panic("not implemented")
}

func (jf joinFault) Level() FaultLevel {
	return jf.level
}

func (jf joinFault) Timestamp() time.Time {
	return jf.timestamp
}

func Join(faults ...Fault) Fault {
	// 1. Remove nil faults.
	faults = RejectNil(faults)
	if len(faults) == 0 {
		return nil
	}

	// 2. Get the highest level of severity.
	highest, ok := HighestFaultLevel(faults)
	if !ok {
		panic("failed to get the highest level of severity")
	}

	// 3. Get the oldest fault.
	oldest, timestamp := OldestFault(faults)
	if oldest == nil {
		panic("failed to get the oldest fault")
	}

	js := &joinFault{
		faults:    faults,
		level:     highest,
		msg:       fmt.Sprintf("joined %d faults", len(faults)),
		timestamp: timestamp,
	}

	return js
}
