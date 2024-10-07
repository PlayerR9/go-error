package faults

import flt "github.com/PlayerR9/go-fault"

type FaultOption func(fault flt.Fault)

func WithAt(at string) FaultOption {
	return func(fault flt.Fault) {
		_ = AddKey(fault, "at", at)
	}
}

func WithBefore(before string) FaultOption {
	return func(fault flt.Fault) {
		_ = AddKey(fault, "before", before)
	}
}

func WithAfter(after string) FaultOption {
	return func(fault flt.Fault) {
		_ = AddKey(fault, "after", after)
	}
}
