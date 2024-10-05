package errs

import "reflect"

var (
	// _FaultType is the reflect type of Fault interface.
	_FaultType reflect.Type
)

func init() {
	type_ := reflect.TypeOf((*Fault)(nil))
	if type_ == nil {
		panic("Fault type is nil")
	}

	_FaultType = type_.Elem()
}

func FilterNonNilFaults(faults []Fault) []Fault {
	if len(faults) == 0 {
		return nil
	}

	var top int

	for i := 0; i < len(faults); i++ {
		if faults[i] != nil {
			faults[top] = faults[i]
			top++
		}
	}

	return faults[:top:top]
}

func Traverse(fault Fault, fn func(fault Fault) bool) bool {
	if fault == nil || fn == nil {
		return false
	}

	stack := []Fault{fault}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if fn(top) {
			return true
		}

		switch fault := top.(type) {
		case interface{ Unwrap() Fault }:
			inner := fault.Unwrap()
			if inner != nil {
				stack = append(stack, inner)
			}
		case interface{ Unwrap() []Fault }:
			inners := fault.Unwrap()

			if len(inners) > 0 {
				var top int

				for i := 0; i < len(inners); i++ {
					if inners[i] != nil {
						inners[top] = inners[i]
						top++
					}
				}

				inners = inners[:top:top]

				stack = append(stack, inners...)
			}
		}
	}

	return false
}
