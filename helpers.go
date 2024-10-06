package errs

import (
	"io"
	"reflect"
)

////////////////////////////////////////////////////////////////////

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

func WriteString(w io.Writer, str string) error {
	if str == "" {
		return nil
	} else if w == nil {
		return io.ErrShortWrite
	}

	data := []byte(str)

	n, err := w.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return io.ErrShortWrite
	}

	return nil
}
