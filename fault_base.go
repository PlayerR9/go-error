package errs

import "time"

type FaultBase[C FaultCode] struct {
	level     FaultLevel
	code      C
	msg       string
	timestamp time.Time
}

func (fb FaultBase[C]) Error() string {
	return "[" + fb.level.String() + "] (" + fb.code.String() + ") " + fb.msg
}

func New[C FaultCode](code C, msg string) Fault {
	return &FaultBase[C]{
		code:      code,
		msg:       msg,
		timestamp: time.Now(),
	}
}
