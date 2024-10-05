package errs

import "io"

func WriteFault(w io.Writer, fault Fault) error {
	if fault == nil {
		return nil
	}

	data := []byte(fault.Error())
	if len(data) == 0 {
		return nil
	} else if w == nil {
		return io.ErrShortWrite
	}

	n, err := w.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return io.ErrShortWrite
	}

	return nil
}
