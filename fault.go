package errs

import "reflect"

// FaultCode is the interface that all error codes should implement.
type FaultCode interface {
	~int

	// String returns the string representation of the error code.
	//
	// Returns:
	//   - string: The string representation of the error code.
	String() string
}

// Fault is the interface that all errors should implement.
type Fault interface {
	// Error returns the string representation of the error.
	//
	// Returns:
	//   - string: The string representation of the error.
	Error() string
}

func Is(fault Fault, target Fault) bool {
	if fault == nil || target == nil {
		return false
	}

	return Traverse(fault, func(fault Fault) bool {
		if fault == target {
			return true
		}

		_, ok := fault.(interface{ Is(Fault) bool })
		return ok
	})
}

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

	return Traverse(fault, func(fault Fault) bool {
		if reflect.TypeOf(fault).AssignableTo(type_) {
			target_value.Elem().Set(reflect.ValueOf(fault))

			return true
		}

		e, ok := fault.(interface{ As(any) bool })
		return ok && e.As(target)
	})
}

func Join(faults ...Fault) Fault {
	panic("not implemented")
}

// Unwrap extracts the underlying fault from a Fault. Yet, it doesn't extract it when the fault implements
// Unwrap() []Fault interface.
//
// Parameters:
//   - fault: The fault to unwrap.
//
// Returns:
//   - Fault: The underlying fault.
//
// This function returns nil when:
// - The fault is nil.
// - The fault implements Unwrap() Fault interface but the underlying fault is nil.
// - The fault doesn't implement Unwrap() Fault interface.
func Unwrap(fault Fault) Fault {
	if fault == nil {
		return nil
	}

	f, ok := fault.(interface{ Unwrap() Fault })
	if !ok {
		return nil
	}

	return f.Unwrap()
}
