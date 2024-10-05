package errs

// Predicate is the type of a function that checks a given condition on a fault.
//
// Parameters:
//   - fault: The fault to check.
//
// Returns:
//   - bool: True if the condition is met, false otherwise.
type Predicate func(fault Fault) bool

// ApplyFilter applies the filter function to the faults and returns the filtered faults.
//
// Parameters:
//   - faults: The faults to apply the filter function to.
//   - fn: The filter function.
//
// Returns:
//   - []Fault: The filtered faults. Nil if no faults are specified or the filter function is nil.
//
// If the fault function returns true, it will be included in the filtered faults.
//
// WARNING:
//   - This function has side-effects; meaning that it will mutate the input faults.
func ApplyFilter(faults []Fault, fn Predicate) []Fault {
	if len(faults) == 0 || fn == nil {
		return nil
	}

	var top int

	for i := 0; i < len(faults); i++ {
		if fn(faults[i]) {
			faults[top] = faults[i]
			top++
		}
	}

	return faults[:top:top]
}

func ApplyReject(faults []Fault, fn Predicate) []Fault {
	if len(faults) == 0 || fn == nil {
		return nil
	}

	var top int

	for i := 0; i < len(faults); i++ {
		if !fn(faults[i]) {
			faults[top] = faults[i]
			top++
		}
	}

	return faults[:top:top]
}
