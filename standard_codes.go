package fault

import "fmt"

// StandardCode is a set of common, standard fault codes.
type StandardCode int

const (
	// Invalid specifies when faults are constructed in a way that is not expected.
	// This is used mostly for internal usage. However, users may use this code
	// when constructing faults in their own code.
	Invalid StandardCode = iota - 1

	// UnknownCode specifies when a fault is created but no code was specified.
	UnknownCode

	// BadParameter specifies a broad category of faults that are caused by unexpected
	// parameters.
	BadParameter

	// OperationFailed specifies a broad category of faults that are returned by functions
	// when they fail.
	OperationFailed
)

var (
	// BadConstruction occurs when a fault does not embed *baseFault.
	BadConstruction FaultDescriber
)

func init() {
	BadConstruction = NewDescriptor(FATAL, Invalid, "fault does not implement *baseFault")
}

// NewBadParameter creates a new BadParameter fault.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func NewBadParameter(msg string, opts ...FaultOption) Fault {
	desc := NewDescriptor(ERROR, BadParameter, msg)

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
//   - Fault: The new Fault. Never returns nil.
func NewNilParameter(param_name string, opts ...FaultOption) Fault {
	desc := NewDescriptor(ERROR, BadParameter, fmt.Sprintf("parameter (%q) must be non-nil", param_name))

	fault := desc.Init()

	for _, opt := range opts {
		opt(fault)
	}

	return fault
}

// NewNilReceiver creates a new OperationFailed fault.
//
// Parameters:
//   - msg: The message of the fault.
//
// Returns:
//   - Fault: The new Fault. Never returns nil.
func NewNilReceiver(opts ...FaultOption) Fault {
	desc := NewDescriptor(ERROR, OperationFailed, "receiver must be non-nil")

	fault := desc.Init()
	_ = SetSuggestions(fault, "Did you forgot to initialize the receiver?")

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
//   - Fault: The new Fault. Never returns nil.
func NewInvalidUsage(message, usage string, opts ...FaultOption) Fault {
	desc := NewDescriptor(ERROR, OperationFailed, message)

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
//   - Fault: The new Fault. Never returns nil.
func NewNoSuchKey(key string, opts ...FaultOption) Fault {
	desc := NewDescriptor(ERROR, OperationFailed, fmt.Sprintf("the specified key (%q) does not exist", key))

	fault := desc.Init()

	for _, opt := range opts {
		opt(fault)
	}

	return fault
}

type FaultOption func(fault Fault)

func WithAt(at string) FaultOption {
	return func(fault Fault) {
		_ = AddKey(fault, "at", at)
	}
}

func WithBefore(before string) FaultOption {
	return func(fault Fault) {
		_ = AddKey(fault, "before", before)
	}
}

func WithAfter(after string) FaultOption {
	return func(fault Fault) {
		_ = AddKey(fault, "after", after)
	}
}
