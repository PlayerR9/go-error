package errs

import (
	"io"
	"strconv"
)

// ErrShortWrite is an error that indicates that the write operation has been
// short written.
type ErrShortWrite struct {
	Fault

	// Written is the number of bytes that have been written.
	Written int

	// Reason is the reason for the short write.
	Reason string
}

// Embeds implements the Fault interface.
//
// Always returns nil.
func (e ErrShortWrite) Embeds() Fault {
	return e.Fault
}

// WriteInfo implements the Fault interface.
//
// Format:
//
//	"Bytes written: <written>"
//	"Reason: <reason>"
//
// Where:
//   - <written> is the number of bytes that have been written.
//   - <reason> is the reason for the short write.
func (e ErrShortWrite) WriteInfo(w Writer) (int, Fault) {
	if w == nil {
		return 0, NewErrShortWrite(0, "no writer provided")
	}

	data := []byte("Bytes written: " + strconv.Itoa(e.Written) + "\nReason: " + e.Reason + "\n")

	n, err := w.Write(data)
	return n, err
}

// NewErrShortWrite creates a new ErrShortWrite.
//
// Parameters:
//   - n: The number of bytes that have been written.
//   - reason: The reason for the short write.
//
// Returns:
//   - *ErrShortWrite: The new ErrShortWrite. Never returns nil.
func NewErrShortWrite(n int, reason string) *ErrShortWrite {
	base := New(OperationFail, "short write")

	fault := &ErrShortWrite{
		Fault:   base,
		Written: n,
		Reason:  reason,
	}

	return fault
}

// Writer is the interface that all writers should implement.
type Writer interface {
	// Write writes the data to the writer.
	//
	// Parameters:
	//   - data: The data to write.
	//
	// Returns:
	//   - int: The number of bytes that have been written.
	//   - *ErrShortWrite: The fault that occurred while writing the data.
	Write(data []byte) (int, *ErrShortWrite)
}

type writer struct {
	w io.Writer
}

func (w writer) Write(data []byte) (int, *ErrShortWrite) {
	n, err := w.w.Write(data)
	if err != nil {
		return n, NewErrShortWrite(n, err.Error())
	} else if n != len(data) {
		return n, NewErrShortWrite(n, "not all bytes written")
	}

	return n, nil
}

func WriterOf(w io.Writer) Writer {
	if w == nil {
		return nil
	}

	return &writer{
		w: w,
	}
}

// Write writes the data to the writer.
//
// Parameters:
//   - w: The writer to write the data to.
//   - data: The data to write.
//
// Returns:
//   - int: The number of bytes that have been written.
//   - *ErrShortWrite: The fault that occurred while writing the data.
//
// If a write is successful, the returned int is guaranteed to be equal to len(data).
func Write(w Writer, data []byte) (int, Fault) {
	if len(data) == 0 {
		return 0, nil
	} else if w == nil {
		return 0, NewErrShortWrite(0, "no writer provided")
	}

	n, err := w.Write(data)
	if err != nil {
		return n, err
	} else if n != len(data) {
		return n, NewErrShortWrite(n, "not all bytes written")
	}

	return n, nil
}

// WriteNewline writes newlines to the writer.
//
// Parameters:
//   - w: The writer to write the newline to.
//   - count: The number of newlines to write. If count is <= 0, no newline is written.
//
// Returns:
//   - int: The number of bytes that have been written.
//   - *ErrShortWrite: The fault that occurred while writing the newline.
//
// If a write is successful, the returned int is guaranteed to be equal to 1.
func WriteNewline(w Writer, count int) (int, Fault) {
	if count <= 0 {
		return 0, nil
	} else if w == nil {
		return 0, NewErrShortWrite(0, "no writer provided")
	}

	data := []byte("\n")

	var total int

	for i := 0; i < count; i++ {
		n, err := w.Write(data)
		total += n

		if err != nil {
			return total, err
		} else if n != len(data) {
			return total, NewErrShortWrite(n, "not all bytes written")
		}
	}

	return total, nil
}

////////////////////////////////////////////////////////////////////

/* func Write(w io.Writer, data []byte) *ErrShortWrite {
	n, err := w.Write(data)
	if err != nil {
		return NewErrShortWrite(n)
	} else if n != len(data) {
		return NewErrShortWrite(n, nil)
	}

	return nil
} */

//////////
/*
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
		return NewFromError(OperationFail, err)
	} else if n != len(data) {
		return ErrShortWrite
	}

	return nil
}
*/
