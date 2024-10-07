package faults

import (
	"fmt"

	flt "github.com/PlayerR9/go-fault"
)

// NewNilReceiver creates a new OperationFailed fault.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - flt.Fault: The new flt.Fault. Never returns nil.
func NewNilReceiver(opts ...FaultOption) flt.Fault {
	desc := flt.NewDescriptor(flt.ERROR, flt.OperationFailed, "receiver must be non-nil")

	fault := desc.Init()
	_ = SetSuggestions(fault, "Did you forgot to initialize the receiver?")

	for _, opt := range opts {
		opt(fault)
	}

	return fault
}

// NewBadParameter creates a new BadParameter fault.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - flt.Fault: The new flt.Fault. Never returns nil.
func NewBadParameter(msg string, opts ...FaultOption) flt.Fault {
	desc := flt.NewDescriptor(flt.ERROR, flt.BadParameter, msg)

	fault := desc.Init()

	for _, opt := range opts {
		opt(fault)
	}

	return fault
}

// NewNilParameter creates a new BadParameter fault.
//
// Parameters:
//   - param_name: The name of the parameter.
//
// Returns:
//   - flt.Fault: The new flt.Fault. Never returns nil.
func NewNilParameter(param_name string, opts ...FaultOption) flt.Fault {
	desc := flt.NewDescriptor(flt.ERROR, flt.BadParameter, fmt.Sprintf("parameter (%q) must be non-nil", param_name))

	fault := desc.Init()

	for _, opt := range opts {
		opt(fault)
	}

	return fault
}

// NewInvalidUsage creates a new OperationFailed fault.
//
// Parameters:
//   - msg: The message of the fault.
//   - usage: The usage of the fault.
//
// Returns:
//   - flt.Fault: The new flt.Fault. Never returns nil.
func NewInvalidUsage(message, usage string, opts ...FaultOption) flt.Fault {
	desc := flt.NewDescriptor(flt.ERROR, flt.OperationFailed, message)

	fault := desc.Init()
	_ = SetSuggestions(fault, usage)

	for _, opt := range opts {
		opt(fault)
	}

	return fault
}

// NewNoSuchKey creates a new OperationFailed fault.
//
// Parameters:
//   - key: The key that was not found.
//
// Returns:
//   - flt.Fault: The new flt.Fault. Never returns nil.
func NewNoSuchKey(key string, opts ...FaultOption) flt.Fault {
	desc := flt.NewDescriptor(flt.ERROR, flt.OperationFailed, fmt.Sprintf("the specified key (%q) does not exist", key))

	fault := desc.Init()

	for _, opt := range opts {
		opt(fault)
	}

	return fault
}
