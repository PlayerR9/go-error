package errs

import (
	"reflect"
	"time"
)

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

// CountNotNil is a helper function that counts the number of non-nil faults in a list of faults.
//
// Parameters:
//   - faults: The list of faults to count.
//
// Returns:
//   - int: The number of non-nil faults in the list.
func CountNotNil(faults []Fault) int {
	var count int

	for _, fault := range faults {
		if fault != nil {
			count++
		}
	}

	return count
}

// RejectNil is a helper function that removes all nil faults from a list of faults.
//
// Parameters:
//   - faults: The list of faults to remove nil faults from.
//
// Returns:
//   - []Fault: The list of faults without nil faults.
func RejectNil(faults []Fault) []Fault {
	var count int

	for _, fault := range faults {
		if fault != nil {
			count++
		}
	}

	if count == 0 {
		return nil
	}

	result := make([]Fault, 0, count)

	for i := 0; i < len(faults); i++ {
		if faults[i] != nil {
			result = append(result, faults[i])
		}
	}

	return result
}

// HighestFaultLevel is a helper function that returns the highest level of severity in a list of
// faults.
//
// Parameters:
//   - faults: The list of faults to get the highest level of severity from.
//
// Returns:
//   - FaultLevel: The highest level of severity in the list.
//   - bool: True if there exists a maximum level of severity in the list, false otherwise.
//
// This function ignores nil faults.
func HighestFaultLevel(faults []Fault) (FaultLevel, bool) {
	faults = RejectNil(faults)
	if len(faults) == 0 {
		return 0, false
	}

	highest := faults[0].Level()

	for _, fault := range faults[1:] {
		level := fault.Level()

		if level > highest {
			highest = level
		}
	}

	return highest, true
}

// OldestFault is a helper function that returns the oldest fault in a list of faults.
//
// Parameters:
//   - faults: The list of faults to get the oldest fault from.
//
// Returns:
//   - Fault: The oldest fault in the list. Nil if there is no oldest fault.
//   - time.Time: The timestamp of the oldest fault in the list.
func OldestFault(faults []Fault) (Fault, time.Time) {
	faults = RejectNil(faults)
	if len(faults) == 0 {
		return nil, time.Time{}
	}

	min := faults[0].Timestamp()
	idx := 0

	for i := 1; i < len(faults); i++ {
		timestamp := faults[i].Timestamp()

		if timestamp.Before(min) {
			min = timestamp
			idx = i
		}
	}

	return faults[idx], min
}
