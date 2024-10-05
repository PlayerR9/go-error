package errs

import "reflect"

type FaultCode interface {
	~int

	String() string
}

type Fault interface {
	Error() string
}

func New[C FaultCode](code C, msg string) Fault {
	return &fault_string[C]{
		code: code,
		msg:  msg,
	}
}

func Is(fault Fault, target Fault) bool {
	if fault == nil || target == nil {
		return false
	}

	return do(fault, func(fault Fault) bool {
		return fault == target
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

	stack := []Fault{fault}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if reflect.TypeOf(top).AssignableTo(type_) {
			target_value.Elem().Set(reflect.ValueOf(top))

			return true
		}

		e, ok := fault.(interface{ As(any) bool })
		if ok && e.As(target) {
			return true
		}

		switch fault := fault.(type) {
		case interface{ Unwrap() Fault }:
			inner := fault.Unwrap()
			if inner != nil {
				stack = append(stack, inner)
			}
		case interface{ Unwrap() []Fault }:
			inners := fault.Unwrap()

			for _, inner := range inners {
				if inner != nil {
					stack = append(stack, inner)
				}
			}
		}
	}

	return false
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
