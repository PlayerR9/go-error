package errs

import "io"

var (
	ErrShortWrite Fault
)

func init() {
	ErrShortWrite = New(OperationFail, "short write")
}

func WriteFault(w io.Writer, fault Fault) Fault {
	if fault == nil {
		return nil
	}

	data := []byte(fault.Error())
	if len(data) == 0 {
		return nil
	} else if w == nil {
		return ErrShortWrite
	}

	n, err := w.Write(data)
	if err != nil {
		return FromError(OperationFail, err)
	} else if n != len(data) {
		return ErrShortWrite
	}

	return nil
}
