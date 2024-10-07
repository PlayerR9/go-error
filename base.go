package fault

import (
	"fmt"
	"slices"
	"strings"
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
	// UnknownLevel represents a fault that has not been initialized or set yet.
	UnknownLevel FaultLevel = iota - 1 // UNKNOWN LEVEL

	// FATAL is the highest level of severity and represents faults that are panic-level of
	// severity.
	FATAL // FATAL

	// ERROR is the second highest level of severity and represents faults that are recoverable
	// errors. This is the "normal" level of severity.
	ERROR // ERROR

	// WARNING is the third highest level of severity and represents faults that are not
	// problematic but may require attention.
	WARNING // WARNING

	// NOTICE is the fourth highest level of severity and are only used to inform the user/
	// operator about something that is not critical.
	NOTICE // NOTICE

	// DEBUG is the lowest level of severity and represents faults that are only used during
	// debugging and ignored in production.
	DEBUG // DEBUG
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

	return &baseFault{
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

// baseFault is the base implementation of the Fault interface.
type baseFault struct {
	// descriptor is the root information of any fault.
	descriptor FaultDescriber

	// timestamp specifies the time when the fault occurred.
	timestamp time.Time

	// suggestions describes one or more possible solutions or actions that can be taken
	// to resolve the fault.
	suggestions []string

	// stack_trace is the stack trace of the fault.
	stack_trace []string

	// context is the context of the fault.
	context map[string]any

	// cause is the cause of the fault.
	cause Fault
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (bf baseFault) Embeds() Fault {
	return nil
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"Occurred at: <timestamp>"
//
// "Suggestions:"
// "- <suggestion>"
// "- ..."
// "Stack trace:"
// "- <stack trace>"
//
// Where:
//   - <timestamp>: The time when the fault occurred.
//   - <suggestion>: One or more possible solutions or actions that can be taken
//     to resolve the fault.
//   - <stack trace>: The stack trace of the fault.
func (bf baseFault) InfoLines() []string {
	var lines []string

	if !bf.timestamp.IsZero() {
		lines = append(lines, "Occurred at: "+bf.timestamp.String())
	}

	if len(bf.suggestions) > 0 {
		lines = append(lines, "Suggestions:")

		for _, suggestion := range bf.suggestions {
			lines = append(lines, "- "+suggestion)
		}
	}

	if len(bf.context) > 0 {
		lines = append(lines, "Context:")

		for k, v := range bf.context {
			lines = append(lines, fmt.Sprintf("- %s: %v", k, v))
		}
	}

	if len(bf.stack_trace) > 0 {
		lines = append(lines, "Stack trace:")

		trace := make([]string, len(bf.stack_trace), len(bf.stack_trace)+1)
		copy(trace, bf.stack_trace)
		trace = append(trace, "")

		slices.Reverse(trace)

		lines = append(lines, "- "+strings.Join(trace, " <- "))
	}

	if bf.cause != nil {
		lines = append(lines, "Caused by:")

		tmp := LinesOf(bf.cause)
		for i := 0; i < len(tmp); i++ {
			tmp[i] = "\t" + tmp[i]
		}

		lines = append(lines, tmp...)
	}

	return lines
}

func DescriptorOf(fault Fault) FaultDescriber {
	if fault == nil {
		return nil
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	return base.descriptor
}

func ErrorOf(fault Fault) string {
	if fault == nil {
		return ""
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	return base.descriptor.String()
}

func LevelOf(fault Fault) FaultLevel {
	if fault == nil {
		return UnknownLevel
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	return base.descriptor.Level()
}

func TimestampOf(fault Fault) time.Time {
	if fault == nil {
		return time.Time{}
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	return base.timestamp
}
