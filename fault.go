// Package errs defines the fault system, which is used to handle errors in a different
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
package errs

import (
	"reflect"
	"slices"
	"time"
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

	// WriteInfo writes the fault's additional information to the writer. Of course, this
	// must exclude both what Error() returns and the embedded fault.
	//
	// Parameters:
	//   - w: The writer to write the information to.
	//
	// Returns:
	//   - int: The number of bytes that have been written.
	//   - Fault: The fault that occurred while writing the information.
	//
	// NOTES:
	// 	- If no faults occurred, the returned int value should be equal to the size of
	// 	the data written to the writer.
	WriteInfo(w Writer) (int, Fault)

	// Error returns the string representation of the fault.
	//
	// Returns:
	//   - string: The string representation of the fault.
	//
	// WARNING:
	// 	- No fault should override this method as it is taken care of the fault's base.
	Error() string

	// Level returns the level of the fault.
	//
	// Returns:
	//   - FaultLevel: The level of the fault.
	//
	// WARNING:
	// 	- No fault should override this method as it is taken care of the fault's base.
	Level() FaultLevel

	// Timestamp returns the timestamp of the fault.
	//
	// Returns:
	//   - time.Time: The timestamp of the fault.
	//
	// WARNING:
	// 	- No fault should override this method as it is taken care of the fault's base.
	Timestamp() time.Time
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

// WriteInfo writes the fault's additional information to the writer by traversing the embedding tower
// of the fault and calling WriteInfo() on each fault. Of course, a newline is written between each
// fault.
//
// Parameters:
//   - w: The writer to write the information to.
//   - fault: The fault whose additional information are to be written.
//
// Returns:
//   - int: The number of bytes that have been written.
//   - Fault: The fault that occurred while writing the information.
//
// This function guarantees that, if no faults occurred, the returned int value should be equal
// to the size of the data written to the writer.
func WriteInfo(w Writer, fault Fault) (int, Fault) {
	if fault == nil {
		return 0, nil
	}

	tower := EmbeddingTower(fault)
	var total int

	// 1. First call.
	elem := tower[0]

	n, err := elem.WriteInfo(w)
	total += n

	if err != nil {
		return total, err
	}

	// 2. Subsequent calls.

	for _, elem := range tower[1:] {
		// Write newline between each fault.
		n, err := WriteNewline(w, 1)
		total += n

		if err != nil {
			return total, err
		}

		// Write info.
		n, err = elem.WriteInfo(w)
		total += n

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

// WriteFault writes both the fault's message and its embedding tower to the writer; adding
// two newlines between them.
//
// Parameters:
//   - w: The writer to write the error information to.
//   - fault: The fault to write.
//
// Returns:
//   - int: The number of bytes that have been written.
//   - Fault: The error that occurred while writing the information.
//
// This function guarantees that, if no faults occurred, the returned int value should be equal
// to the size of the data written to the writer.
func WriteFault(w Writer, fault Fault) (int, Fault) {
	if fault == nil {
		return 0, nil
	}

	var total int

	data := []byte(fault.Error())
	if len(data) > 0 {
		n, err := Write(w, data)
		total += n

		if err != nil {
			return total, err
		}

		n, err = WriteNewline(w, 2)
		total += n

		if err != nil {
			return total, err
		}
	}

	n, err := WriteInfo(w, fault)
	total += n

	if err != nil {
		return total, err
	}

	return total, nil
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

	ok := Traverse(fault, func(f Fault) bool {
		if f == target {
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

// try is a helper function for Try.
//
// Parameters:
//   - fault: The fault to recover from.
//   - fn: The function to execute.
//
// Assertions:
//   - fn must not be nil.
//   - fault must not be nil.
func try(fault *Fault, fn func()) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		switch r := r.(type) {
		case Fault:
			*fault = r
		case string:
			*fault = NewErrPanic(r)
		case error:
			*fault = NewFaultErr(r)
		default:
			*fault = NewErrPanic(r)
		}
	}()

	fn()
}

// Try executes a panicing function and returns a fault with the paniced value.
//
// Parameters:
//   - fn: The function to execute.
//
// Returns:
//   - Fault: The fault with the paniced value. Nil if no panic occurred.
//
// Behaviors:
//   - If the panic value is nil or it does not panic, it returns nil.
//   - If the panic value is Fault, it returns it.
//   - If the panic value is error, it returns a new FaultErr with the error.
//   - In all other cases, it returns a new ErrPanic with the panic value.
func Try(fn func()) Fault {
	if fn == nil {
		return nil
	}

	var fault Fault

	try(&fault, fn)

	return fault
}
