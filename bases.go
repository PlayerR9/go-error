package fault

import (
	"strings"
	"time"
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
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(bf.level.String())
	builder.WriteString("] (")
	builder.WriteString(bf.code.String())
	builder.WriteString(") ")
	builder.WriteString(bf.msg)

	return builder.String()
}

// Level implements the Fault interface.
func (bf baseFault[C]) Level() FaultLevel {
	return bf.level
}

// Timestamp implements the Fault interface.
func (bf baseFault[C]) Timestamp() time.Time {
	return bf.timestamp
}
