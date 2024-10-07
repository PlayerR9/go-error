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
			*fault = FromErr(r)
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
//   - fault.Fault: The fault with the paniced value. Nil if no panic occurred.
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
