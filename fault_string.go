package errs

type fault_string[C FaultCode] struct {
	code C
	msg  string
}

func (f fault_string[C]) Error() string {
	return "(" + f.code.String() + ") " + f.msg
}
