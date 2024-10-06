package errs

import (
	"time"
)

// FaultBase is the base implementation of the Fault interface. Every fault must
// embed this struct in order to implement the Fault interface.
type FaultBase[C FaultCode] struct {
	// level is the level of the fault.
	level FaultLevel

	// code is the code of the fault.
	code C

	// msg is the message of the fault.
	msg string

	// timestamp is the timestamp of the fault.
	timestamp time.Time
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (fb FaultBase[C]) Embeds() Fault {
	return nil
}

// Error implements the Fault interface.
//
// Format:
//
//	"[<level>] (<code>) <msg>"
//
// where:
//   - <level>: The level of the fault.
//   - <code>: The code of the fault.
//   - <msg>: The message of the fault.
func (fb FaultBase[C]) Error() string {
	return "[" + fb.level.String() + "] (" + fb.code.String() + ") " + fb.msg
}

// WriteInfo implements the Fault interface.
//
// Format:
//
//	"Occurred at: <timestamp>"
//
// Where, <timestamp> is the timestamp of the fault.
func (fb FaultBase[C]) WriteInfo(w Writer) (int, Fault) {
	var total int

	if !fb.timestamp.IsZero() {
		if w == nil {
			return total, NewErrShortWrite(0, "no writer provided")
		}

		data := []byte("Occurred at: " + fb.timestamp.String())

		n, err := Write(w, data)
		total += n

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

// New creates a new FaultBase struct with the given code and message. The timestamp
// is set to the current time.
//
// Parameters:
//   - code: The code of the fault.
//   - msg: The message of the fault.
//
// Returns:
//   - *FaultBase: The new FaultBase struct. Never returns nil.
func New[C FaultCode](code C, msg string) *FaultBase[C] {
	return &FaultBase[C]{
		level:     ERROR,
		code:      code,
		msg:       msg,
		timestamp: time.Now(),
	}
}

// NewWithLevel creates a new FaultBase struct with the given level, code and message.
// The timestamp is set to the current time.
//
// Parameters:
//   - level: The level of the fault.
//   - code: The code of the fault.
//   - msg: The message of the fault.
//
// Returns:
//   - *FaultBase: The new FaultBase struct. Never returns nil.
func NewWithLevel[C FaultCode](level FaultLevel, code C, msg string) *FaultBase[C] {
	return &FaultBase[C]{
		level:     level,
		code:      code,
		msg:       msg,
		timestamp: time.Now(),
	}
}

/////////////////////////////////

func NewFromError[C FaultCode](code C, err error) *FaultBase[C] {
	var msg string

	if err != nil {
		msg = err.Error()
	}

	return New(code, msg)
}
