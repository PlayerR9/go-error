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
	"reflect"
	"slices"
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

// LinesOf returns the fault's message and its embedding tower as a list of strings.
//
// Parameters:
//   - fault: The fault whose additional information are to be written.
//
// Returns:
//   - []string: The fault's additional information.
//
// An empty line is added between the message and the embedding tower. Also, a "dot" is
// added at the end of the message.
func LinesOf(fault Fault) []string {
	if fault == nil {
		return nil
	}

	var lines []string

	lines = append(lines, ErrorOf(fault)+".")
	lines = append(lines, "")

	tmp := InfoLines(fault)
	lines = append(lines, tmp...)

	return lines
}

// Is checks whether a fault or any fault it may wrap is equal to the target fault. The
// search is done in a depth-first manner and it doesn't scan the fault's embedding tower.
//
// The comparison is done using pointer equality and the Is() method of the target fault. When
// either is true, the function returns true.
//
// Parameters:
//   - fault: The fault to check.
//   - target: The fault to compare with.
//
// Returns:
//   - bool: True if the fault (or any fault it may wrap) is equal to the target fault, false
//     otherwise.
func Is(fault Fault, target Fault) bool {
	if fault == nil || target == nil {
		return false
	}

	target_desc := DescriptorOf(target)

	ok := Traverse(fault, func(f Fault) bool {
		f_desc := DescriptorOf(f)

		if f == target || f_desc == target_desc {
			return true
		}

		_, ok := f.(interface{ Is(Fault) bool })
		return ok
	})

	return ok
}

// As checks whether a fault or any fault it may wrap implements the target interface. The
// search is done in a depth-first manner and it doesn't scan the fault's embedding tower.
//
// The comparison is done using pointer equality and the As() method of the target interface. When
// either is true, the function returns true and the target is set to the fault.
//
// Parameters:
//   - fault: The fault to check.
//   - target: The target to set if the fault implements the target interface.
//
// Returns:
//   - bool: True if the fault (or any fault it may wrap) implements the target interface, false
//     otherwise.
func As(fault Fault, target any) bool {
	if fault == nil || target == nil {
		return false
	}

	target_value := reflect.ValueOf(target)

	target_type := target_value.Type()
	if target_type.Kind() != reflect.Ptr || target_value.IsNil() {
		return false
	}

	type_ := target_type.Elem()
	if type_.Kind() != reflect.Interface && !type_.Implements(_FaultType) {
		return false
	}

	return Traverse(fault, func(f Fault) bool {
		if reflect.TypeOf(f).AssignableTo(type_) {
			target_value.Elem().Set(reflect.ValueOf(f))

			return true
		}

		e, ok := f.(interface{ As(any) bool })
		return ok && e.As(target)
	})
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
