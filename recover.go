package errs

import "fmt"

////////////////////////////////////////////////////////////////////

func try(fault *Fault, fn func()) {
	if fn == nil {
		return
	}

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		switch r := r.(type) {
		case Fault:
			*fault = r
		case string:
		case interface{ Error() string }:
			New()
		default:
			fmt.Printf("panic: %v\n", r)
		}
	}()

	fn()
}

func Try(fn func()) Fault {
	if fn == nil {
		return nil
	}

	var fault Fault

	try(&fault, fn)

	return fault
}
