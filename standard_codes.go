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

	// FaultJoin is the fault code for when a fault is joined.
	FaultJoin

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

// ErrPanic is an error that indicates that a panic occurred.
type ErrPanic struct {
	Fault

	// Value is the value that was passed to panic.
	Value any
}

// Embeds implements the Fault interface.
func (e ErrPanic) Embeds() Fault {
	return e.Fault
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Value: <value>"
func (e ErrPanic) InfoLines() []string {
	lines := make([]string, 0, 1)

	lines = append(lines, "- Value: "+fmt.Sprintf("%v", e.Value))

	return lines
}

// NewErrPanic creates a new ErrPanic.
//
// Parameters:
//   - value: The value that was passed to panic.
//
// Returns:
//   - *ErrPanic: A new ErrPanic. Never returns nil.
func NewErrPanic(value any) *ErrPanic {
	base := WithLevel(FATAL, UnknownCode, "a panic occurred")

	return &ErrPanic{
		Fault: base,
		Value: value,
	}
}

// ErrFault is a fault that wraps a Go error.
type ErrFault struct {
	Fault

	// Err is the error that occurred.
	Err error
}

// Embeds implements the Fault interface.
func (e ErrFault) Embeds() Fault {
	return e.Fault
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Error: <error>"
//
// Where, <error> is the error that occurred. If no error was provided, "no error
// provided" is used instead.
func (e ErrFault) InfoLines() []string {
	lines := make([]string, 0, 1)

	if e.Err != nil {
		lines = append(lines, "- Error: "+e.Err.Error())
	} else {
		lines = append(lines, "- Error: no error provided")
	}

	return lines
}

// FromErr creates a new ErrFault.
//
// Parameters:
//   - err: The error that occurred.
//
// Returns:
//   - *ErrFault: The new ErrFault. Never returns nil.
func FromErr(err error) Fault {
	base := New(UnknownCode, "something went wrong")

	return &ErrFault{
		Fault: base,
		Err:   err,
	}
}

// Join is a helper function that joins a list of faults into a single fault.
type JoinFault struct {
	Fault

	// faults contains all the faults that have been joined.
	faults []Fault
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (jf JoinFault) Embeds() Fault {
	return jf.Fault
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Faults: <faults>"
func (jf JoinFault) InfoLines() []string {
	var lines []string

	for _, fault := range jf.faults {
		tmp := fault.InfoLines()
		lines = append(lines, tmp...)
	}

	return lines
}

// Join is a helper function that joins a list of faults into a single fault.
//
// Parameters:
//   - faults: The faults to join. May be nil.
//
// Returns:
//   - Fault: The joined fault.
//
// This function returns nil if all the faults are nil.
func Join(faults ...Fault) Fault {
	// 1. Remove nil faults.
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

	// 2. Get the highest level of severity.
	highest := LevelOf(faults[0])

	for _, fault := range faults[1:] {
		level := LevelOf(fault)

		if level > highest {
			highest = level
		}
	}

	base := WithLevel(highest, FaultJoin, fmt.Sprintf("joined %d faults", len(faults)))

	js := &JoinFault{
		Fault:  base,
		faults: faults,
	}

	return js
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
