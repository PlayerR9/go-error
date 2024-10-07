package fault

// get_base is a helper function to get the baseFault from a fault.
//
// Parameters:
//   - fault: The fault to get the baseFault from.
//
// Returns:
//   - *baseFault: The baseFault. Returns nil if the fault is nil or does not implement the baseFault interface.
func get_base(fault Fault) *baseFault {
	for fault != nil {
		base, ok := fault.(*baseFault)
		if ok {
			return base
		}

		fault = fault.Embeds()
	}

	return nil
}

// SetSuggestions sets the fault's suggestions; ignoring any empty suggestions.
//
// Parameters:
//   - fault: The fault to set the suggestions for.
//   - suggestions: The suggestions to set.
//
// Returns:
//   - bool: True if the suggestions were set, false otherwise. It only returns false
//     when there is at least one non-empty suggestion and the fault is nil.
func SetSuggestions(fault Fault, suggestions ...string) bool {
	var count int

	for i := 0; i < len(suggestions); i++ {
		if suggestions[i] != "" {
			count++
		}
	}

	if count == 0 {
		return true
	}

	if fault == nil {
		return false
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	filtered := make([]string, 0, count)

	for i := 0; i < len(suggestions); i++ {
		if suggestions[i] != "" {
			filtered = append(filtered, suggestions[i])
		}
	}

	base.suggestions = append(base.suggestions, filtered...)

	return true
}

// Throw adds a stack trace's frame to the fault and returns the fault.
//
// Parameters:
//   - frame: The stack trace's frame to add.
//
// Returns:
//   - Fault: The fault. Returns nil if the fault is nil.
func Throw(fault Fault, frame string) Fault {
	if fault == nil {
		return nil
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	base.stack_trace = append(base.stack_trace, frame)

	return fault
}
