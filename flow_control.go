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

func SetSuggestions(fault Fault, suggestions ...string) {
	var count int

	for i := 0; i < len(suggestions); i++ {
		if suggestions[i] != "" {
			count++
		}
	}

	if count == 0 {
		return
	}

	if fault == nil {
		panic("TODO: Handle nil fault")
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	for i := 0; i < len(suggestions); i++ {
		if suggestions[i] != "" {
			base.suggestions = append(base.suggestions, suggestions[i])
		}
	}
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
		panic("TODO: Handle nil fault")
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	base.stack_trace = append(base.stack_trace, frame)

	return fault
}
