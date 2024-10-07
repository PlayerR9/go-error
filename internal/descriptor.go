package internal

import "strings"

// FaultCode is the interface that all fault codes should implement.
type FaultCode interface {
	~int

	// String returns the string representation of the fault code.
	//
	// Returns:
	//   - string: The string representation of the fault code.
	String() string
}

// Descriptor is the root information of any fault.
type Descriptor[C FaultCode] struct {
	// level indicates the severity level of the fault.
	level FaultLevel

	// code specifies the broader category that the fault belongs to.
	code C

	// msg informs about the nature of the fault.
	msg string
}

// Error returns the string representation of the fault.
//
// Returns:
//   - string: The string representation of the fault.
func (d Descriptor[C]) Error() string {
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(d.level.String())
	builder.WriteString("] (")
	builder.WriteString(d.code.String())
	builder.WriteString(") ")
	builder.WriteString(d.msg)

	return builder.String()
}
