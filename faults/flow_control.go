package faults

import flt "github.com/PlayerR9/go-fault"

// Throw adds a stack trace's frame to the fault and returns the fault.
//
// Parameters:
//   - frame: The stack trace's frame to add.
//
// Returns:
//   - flt.Fault: The fault. Returns nil if the fault is nil.
func Throw(fault flt.Fault, frame string) flt.Fault {
	if fault == nil {
		return nil
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	base.AppendFrame(frame)

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
func try(fault *flt.Fault, fn func()) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		switch r := r.(type) {
		case flt.Fault:
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
//   - fault.flt.Fault: The fault with the paniced value. Nil if no panic occurred.
//
// Behaviors:
//   - If the panic value is nil or it does not panic, it returns nil.
//   - If the panic value is flt.Fault, it returns it.
//   - If the panic value is error, it returns a new FaultErr with the error.
//   - In all other cases, it returns a new ErrPanic with the panic value.
func Try(fn func()) flt.Fault {
	if fn == nil {
		return nil
	}

	var fault flt.Fault

	try(&fault, fn)

	return fault
}
