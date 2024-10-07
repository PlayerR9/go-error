package faults

import (
	"reflect"

	flt "github.com/PlayerR9/go-fault"
)

var (
	// _FaultType is the reflect type of Fault interface.
	_FaultType reflect.Type
)

func init() {
	type_ := reflect.TypeOf((*flt.Fault)(nil))
	if type_ == nil {
		panic("Fault type is nil")
	}

	_FaultType = type_.Elem()
}

// Traverse traverses the fault tree in a DFS manner and executes the function on each fault;
// stopping at the first time the function returns true.
//
// Parameters:
//   - fault: The fault to traverse.
//   - fn: The function to execute on each fault.
//
// Returns:
//   - bool: True if the function returns true, false otherwise.
//
// Behaviors:
//   - If the fault is nil or the function is nil, the function returns false.
func Traverse(fault flt.Fault, fn func(fault flt.Fault) bool) bool {
	if fault == nil || fn == nil {
		return false
	}

	stack := []flt.Fault{fault}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if fn(top) {
			return true
		}

		switch fault := top.(type) {
		case interface{ Unwrap() flt.Fault }:
			inner := fault.Unwrap()
			if inner != nil {
				stack = append(stack, inner)
			}
		case interface{ Unwrap() []flt.Fault }:
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
func Is(fault flt.Fault, target flt.Fault) bool {
	if fault == nil || target == nil {
		return false
	}

	target_desc := DescriptorOf(target)

	ok := Traverse(fault, func(f flt.Fault) bool {
		f_desc := DescriptorOf(f)

		if f == target || f_desc == target_desc {
			return true
		}

		_, ok := f.(interface{ Is(flt.Fault) bool })
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
func As(fault flt.Fault, target any) bool {
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

	return Traverse(fault, func(f flt.Fault) bool {
		if reflect.TypeOf(f).AssignableTo(type_) {
			target_value.Elem().Set(reflect.ValueOf(f))

			return true
		}

		e, ok := f.(interface{ As(any) bool })
		return ok && e.As(target)
	})
}
