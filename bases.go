package fault

import "time"

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

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Occurred at: <timestamp>"
//
// Where, <timestamp> is the timestamp of the fault.
func (bf baseFault[C]) InfoLines() []string {
	var lines []string

	if !bf.timestamp.IsZero() {
		lines = append(lines, "- Occurred at: "+bf.timestamp.String())
	}

	return lines
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
