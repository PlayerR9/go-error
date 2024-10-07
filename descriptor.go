package fault

// FaultCode is the interface that all fault codes should implement.
type FaultCode interface {
	~int

	// String returns the string representation of the fault code.
	//
	// Returns:
	//   - string: The string representation of the fault code.
	String() string
}

// FaultLevel is the severity level of a fault.
type FaultLevel int

const (
	// FATAL is the highest level of severity and represents faults that are panic-level of
	// severity.
	FATAL FaultLevel = iota

	// ERROR is the second highest level of severity and represents faults that are recoverable
	// errors. This is the "normal" level of severity.
	ERROR

	// WARNING is the third highest level of severity and represents faults that are not
	// problematic but may require attention.
	WARNING

	// NOTICE is the fourth highest level of severity and are only used to inform the user/
	// operator about something that is not critical.
	NOTICE

	// DEBUG is the lowest level of severity and represents faults that are only used during
	// debugging and ignored in production.
	DEBUG
)
