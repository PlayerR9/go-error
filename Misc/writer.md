package errs

import (
	"fmt"
	"io"
)

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
	Write(data []byte) (int, Fault)
}

// ErrShortWrite is an error that indicates that the write operation has been
// short written.
type ErrShortWrite struct {
	Fault

	// Expected is the amount of bytes that were expected to be written.
	Expected int

	// Actual is the actual number of bytes that have been written.
	Actual int
}

// Embeds implements the Fault interface.
func (e ErrShortWrite) Embeds() Fault {
	return e.Fault
}

// InfoLines implements the Fault interface.
func (e ErrShortWrite) InfoLines(w Writer) (int, Fault) {
	// - Expected: <expected>
	// - Actual: <actual>
	data := []byte(fmt.Sprintf(
		"- Expected: %d\n- Actual: %d",
		e.Expected,
		e.Actual,
	))

	n, err := Write(w, data)
	return n, err
}

// NewErrShortWrite creates a new ErrShortWrite.
//
// Parameters:
//   - expected: The amount of bytes that were expected to be written.
//   - actual: The actual number of bytes that have been written.
//
// Returns:
//   - *ErrShortWrite: The new ErrShortWrite. Never returns nil.
func NewErrShortWrite(expected, actual int) *ErrShortWrite {
	base := New(OperationFail, "short write")

	err := &ErrShortWrite{
		Fault: base,

		Expected: 0,
		Actual:   0,
	}

	return err
}

// StreamWriter is an adapter that allows interoperation with the io.Writer
// interface.
type StreamWriter struct {
	// w is the underlying io.Writer.
	w io.Writer

	// err is the error that occurred while writing.
	err error

	// actual is the number of bytes that have been written.
	actual int

	// total is the total number of bytes that have been written.
	total int
}

// Write implements the io.Writer interface.
//
// Panics if the receiver is nil.
func (sw *StreamWriter) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	} else if sw == nil {
		panic("receiver must not be nil")
	}

	sw.total += len(data)

	n, err := sw.w.Write(data)
	sw.actual += n

	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	if err != nil {
		sw.err = err
	}

	return n, err
}

// Fault returns the fault that occurred while writing.
//
// Returns:
//   - Fault: The fault that occurred while writing. If there was no fault, it will be nil.
//
// Faults:
//   - ErrShortWrite: The fault that occurred while writing.
//   - FaultErr: In any other case.
func (sw StreamWriter) Fault() Fault {
	if sw.err != nil {
		return NewFaultErr(sw.err)
	}

	if sw.actual < sw.total {
		return NewErrShortWrite(sw.total, sw.actual)
	}

	return nil
}

// WriterOf creates a StreamWriter from an io.Writer.
//
// Parameters:
//   - w: The io.Writer to create the StreamWriter from. If w is nil, io.Discard is used.
//
// Returns:
//   - *StreamWriter: The new StreamWriter. Never returns nil.
func WriterOf(w io.Writer) *StreamWriter {
	if w == nil {
		w = io.Discard
	}

	return &StreamWriter{
		w: w,
	}
}

///////////////////////////////////////

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
	panic("not implemented")
	// if len(data) == 0 {
	// 	return 0, nil
	// } else if w == nil {
	// 	return 0, NewErrShortWrite(0, "no writer provided")
	// }

	// n, err := w.Write(data)
	// if err != nil {
	// 	return n, err
	// } else if n != len(data) {
	// 	return n, NewErrShortWrite(n, "not all bytes written")
	// }

	// return n, nil
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
	panic("not implemented")

	// if count <= 0 {
	// 	return 0, nil
	// } else if w == nil {
	// 	return 0, NewErrShortWrite(0, "no writer provided")
	// }

	// data := []byte("\n")

	// var total int

	// for i := 0; i < count; i++ {
	// 	n, err := w.Write(data)
	// 	total += n

	// 	if err != nil {
	// 		return total, err
	// 	} else if n != len(data) {
	// 		return total, NewErrShortWrite(n, "not all bytes written")
	// 	}
	// }

	// return total, nil
}
