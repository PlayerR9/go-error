package errs

import (
	"fmt"
	"time"
)

// FaultCode is the interface that all fault codes should implement.
type FaultCode interface {
	~int

	// String returns the string representation of the fault code.
	//
	// Returns:
	//   - string: The string representation of the fault code.
	String() string
}

// FaultLevel is the severity level of a fault.
type FaultLevel int

const (
	// FATAL is the highest level of severity and represents faults that are panic-level of
	// severity.
	FATAL FaultLevel = iota

	// ERROR is the second highest level of severity and represents faults that are recoverable
	// errors. This is the "normal" level of severity.
	ERROR

	// WARNING is the third highest level of severity and represents faults that are not
	// problematic but may require attention.
	WARNING

	// NOTICE is the fourth highest level of severity and are only used to inform the user/
	// operator about something that is not critical.
	NOTICE

	// DEBUG is the lowest level of severity and represents faults that are only used during
	// debugging and ignored in production.
	DEBUG
)

// baseFault is the base implementation of the Fault interface.
type baseFault[C FaultCode] struct {
	// level indicates the severity level of the fault.
	level FaultLevel

	// code specifies the broader category that the fault belongs to.
	code C

	// msg informs about the nature of the fault.
	msg string

	// timestamp specifies the time when the fault occurred.
	timestamp time.Time
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (bf baseFault[C]) Embeds() Fault {
	return nil
}

// WriteInfo implements the Fault interface.
//
// Format:
//
//	"- Occurred at: <timestamp>"
//
// Where, <timestamp> is the timestamp of the fault.
func (bf baseFault[C]) WriteInfo(w Writer) (int, Fault) {
	var total int

	if !bf.timestamp.IsZero() {
		data := []byte("Occurred at: " + bf.timestamp.String())

		n, err := Write(w, data)
		total += n

		if err != nil {
			return total, err
		}
	}

	return total, nil
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
func (bf baseFault[C]) Error() string {
	return "[" + bf.level.String() + "] (" + bf.code.String() + ") " + bf.msg
}

// Level implements the Fault interface.
func (bf baseFault[C]) Level() FaultLevel {
	return bf.level
}

// Timestamp implements the Fault interface.
func (bf baseFault[C]) Timestamp() time.Time {
	return bf.timestamp
}

// FaultErr is a fault that wraps a Go error.
type FaultErr struct {
	Fault

	// Err is the error that occurred.
	Err error
}

// Embeds implements the Fault interface.
func (e FaultErr) Embeds() Fault {
	return e.Fault
}

// WriteInfo implements the Fault interface.
//
// Format:
//
//	"- Error: <error>"
//
// Where, <error> is the error that occurred. If no error was provided, "no error
// provided" is used instead.
func (e FaultErr) WriteInfo(w Writer) (int, Fault) {
	var err_str string

	if e.Err == nil {
		err_str = "no error provided"
	} else {
		err_str = e.Err.Error()
	}

	data := []byte("- Error: " + err_str)

	n, err := Write(w, data)

	return n, err
}

// NewFaultErr creates a new FaultErr.
//
// Parameters:
//   - err: The error that occurred.
//
// Returns:
//   - *FaultErr: The new FaultErr. Never returns nil.
func NewFaultErr(err error) Fault {
	base := New(Unknown, "something went wrong")

	return &FaultErr{
		Fault: base,
		Err:   err,
	}
}

// ErrPanic is an error that indicates that a panic occurred.
type ErrPanic struct {
	Fault

	// Value is the value that was passed to panic.
	Value any
}

// Embeds implements the Fault interface.
func (e ErrPanic) Embeds() Fault {
	return e.Fault
}

// WriteInfo implements the Fault interface.
//
// Format:
//
//	"- Value: <value>"
func (e ErrPanic) WriteInfo(w Writer) (int, Fault) {
	data := []byte(fmt.Sprintf("- Value: %v", e.Value))

	n, err := Write(w, data)
	return n, err
}

// NewErrPanic creates a new ErrPanic.
//
// Parameters:
//   - value: The value that was passed to panic.
//
// Returns:
//   - *ErrPanic: A new ErrPanic. Never returns nil.
func NewErrPanic(value any) *ErrPanic {
	base := WithLevel(FATAL, Unknown, "a panic occurred")

	return &ErrPanic{
		Fault: base,
		Value: value,
	}
}
