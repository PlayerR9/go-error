// Package fault defines the fault system, which is used to handle errors in a different
// way than the standard library.
//
// A fault, unlike with the Go error, is defined to be a more "complex" data type that,
// aside from having a message, it also carries additional information that follows the
// "standard" format of:
//
//	"[<level>] (<code>) <msg>"
//
//	Occurred at: <timestamp>
//
//	// Additional information...
//
// where:
//   - <level>: The level of the fault.
//   - <code>: The code of the fault.
//   - <msg>: The message of the fault.
//   - <timestamp>: The timestamp of the fault.
//   - Additional information: The additional information of the fault.
//
// Generally speaking, a custom fault is declared in the following way:
//
//	type ErrMyFault struct {
//		Fault
//
//		// Additional fields...
//	}
//
//	func (e ErrMyFault) Embeds() Fault {
//		return e.Fault
//	}
//
//	func (e ErrMyFault) WriteInfo(w Writer) (int, Fault) {
//		// Write here the additional information of the fault. (Do not call e.Fault.WriteInfo()!)
//	}
//
// Here, ErrMyFault embeds another fault and implements the Fault interface. As you can see,
// the fault does not implement the Error() method of the Fault interface as, it's up to the
// embedded fault to implement it.
//
// Therefore, any constructor would look like:
//
//	func NewErrMyFault(/*<args>*/) *ErrMyFault {
//		base := // Constructor of the base fault
//
//		return &ErrMyFault{
//			Fault: base,
//
//			// Additional fields...
//		}
//	}
package fault

import (
	"slices"
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

// Fault is implemented by all errors/faults. However, a fault must embed another fault in order to implement
// this interface. The embedded fault is referred to as the "base".
type Fault interface {
	// Embeds returns the base of the fault.
	//
	// Returns:
	//   - Fault: The base of the fault.
	//
	// A return value of nil indicates either that the fault is a FaultBase or that
	// the fault does not want to "show" the embedded value.
	Embeds() Fault

	// InfoLines returns the fault's additional information as a list of strings. The result
	// is anything that is not already specified by the Error() method and its embedding tower.
	//
	// Returns:
	//   - []string: The fault's additional information.
	InfoLines() []string
}

// EmbeddingTower returns a list of all the bases that make up the embedding tower of the fault,
// from the innermost fault to the outermost fault.
//
// An embedding tower is obtained by calling Embeds() on the fault until either the fault
// does not embed any other fault or the fault was already seen. (The latter case prevents
// infinite loops.)
//
// Parameters:
//   - fault: The fault to get the embedding tower of.
//
// Returns:
//   - []Fault: The embedding tower of the fault. The first element is the innermost fault
//     and the last element is the outermost fault.
func EmbeddingTower(fault Fault) []Fault {
	if fault == nil {
		return nil
	}

	seen := make(map[Fault]struct{})

	stack := []Fault{fault}

	for {
		top := stack[len(stack)-1]

		_, ok := seen[top]
		if ok {
			break
		}

		seen[top] = struct{}{}

		embedded := top.Embeds()
		if embedded == nil {
			break
		}

		stack = append(stack, embedded)
	}

	slices.Reverse(stack)

	return stack
}

// InfoLines returns the fault's additional information as a list of strings. The result is obtained
// by traversing the embedding tower of the fault and calling InfoLines() on each fault.
//
// Parameters:
//   - fault: The fault whose additional information are to be written.
//
// Returns:
//   - []string: The fault's additional information.
func InfoLines(fault Fault) []string {
	if fault == nil {
		return nil
	}

	tower := EmbeddingTower(fault)

	var lines []string

	for _, elem := range tower {
		tmp := elem.InfoLines()
		lines = append(lines, tmp...)
	}

	return lines
}

// Unwrap extracts the underlying fault from a Fault. It doesn't affect the embedding tower and ignores the error
// if it does not implement Unwrap() Fault method.
//
// Parameters:
//   - fault: The fault to unwrap.
//
// Returns:
//   - Fault: The underlying fault.
func Unwrap(fault Fault) Fault {
	if fault == nil {
		return nil
	}

	uw, ok := fault.(interface{ Unwrap() Fault })
	if !ok {
		return nil
	}

	return uw.Unwrap()
}
