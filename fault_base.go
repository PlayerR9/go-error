package errs

import (
	"io"
	"time"
)

type FaultBase[C FaultCode] struct {
	level     FaultLevel
	code      C
	msg       string
	timestamp time.Time
}

func (fb FaultBase[C]) Error() string {
	return "[" + fb.level.String() + "] (" + fb.code.String() + ") " + fb.msg
}

func (fb FaultBase[C]) WriteInfo(w io.Writer) Fault {
	if !fb.timestamp.IsZero() {
		// ("Occurred at: %v\n", fb.timestamp)
	}

	return nil
}

func New[C FaultCode](code C, msg string) Fault {
	return &FaultBase[C]{
		level:     ERROR,
		code:      code,
		msg:       msg,
		timestamp: time.Now(),
	}
}

func NewWithLevel[C FaultCode](level FaultLevel, code C, msg string) Fault {
	return &FaultBase[C]{
		level:     level,
		code:      code,
		msg:       msg,
		timestamp: time.Now(),
	}
}
