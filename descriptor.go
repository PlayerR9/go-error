package fault

import (
	"strings"
	"time"
)

// FaultDescriber is the interface that all fault describers should implement.
type FaultDescriber interface {
	// String returns the error message of the fault.
	//
	// Returns:
	//   - string: The error message of the fault.
	String() string

	// Level returns the severity level of the fault.
	//
	// Returns:
	//   - FaultLevel: The severity level of the fault.
	Level() FaultLevel

	// Init initializes the fault describer by creating a new Fault instance.
	//
	// Returns:
	//   - Fault: The new Fault. Never returns nil, unless the receiver is nil.
	Init() Fault
}

// faultDescriptor is the root information of any fault. Once created, it is read-only.
type faultDescriptor[C FaultCode] struct {
	// level indicates the severity level of the fault.
	level FaultLevel

	// code specifies the broader category that the fault belongs to.
	code C

	// msg informs about the nature of the fault.
	msg string
}

// String implements the fmt.Stringer interface.
//
// Format:
//
//	"[<level>] (<code>) <msg>"
//
// where:
//   - <level>: The level of the fault.
//   - <code>: The code of the fault.
//   - <msg>: The message of the fault.
func (fd faultDescriptor[C]) String() string {
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(fd.level.String())
	builder.WriteString("] (")
	builder.WriteString(fd.code.String())
	builder.WriteString(") ")
	builder.WriteString(fd.msg)

	return builder.String()
}

// Level implements the Fault interface.
func (fd faultDescriptor[C]) Level() FaultLevel {
	return fd.level
}

// Init implements the FaultDescriber interface.
func (fd *faultDescriptor[C]) Init() Fault {
	if fd == nil {
		return nil
	}

	return &BaseFault{
		descriptor: fd,
		timestamp:  time.Now(),
	}
}

// NewDescriptor creates a new FaultDescriber instance. Each descriptor is unique and
// read-only. As such, comparation can only be done with pointer equality.
//
// Parameters:
//   - level: The level of the fault.
//   - code: The code of the fault.
//   - msg: The message of the fault.
//
// Returns:
//   - FaultDescriber: The new FaultDescriber. Never returns nil.
func NewDescriptor[C FaultCode](level FaultLevel, code C, msg string) FaultDescriber {
	return &faultDescriptor[C]{
		level: level,
		code:  code,
		msg:   msg,
	}
}
