package faults

import (
	"fmt"
	"time"

	flt "github.com/PlayerR9/go-fault"
)

func Access[T flt.Fault](fault flt.Fault) (T, bool) {
	zero := *new(T)

	for fault != nil {
		v, ok := fault.(T)
		if ok {
			return v, true
		}

		fault = fault.Embeds()
	}

	return zero, false
}

func DescriptorOf(fault flt.Fault) flt.FaultDescriber {
	if fault == nil {
		return nil
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	return base.descriptor
}

func ErrorOf(fault flt.Fault) string {
	if fault == nil {
		return ""
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	return base.Error()
}

func LevelOf(fault flt.Fault) flt.FaultLevel {
	if fault == nil {
		return flt.UnknownLevel
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	return base.descriptor.Level()
}

func TimestampOf(fault flt.Fault) time.Time {
	if fault == nil {
		return time.Time{}
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	return base.timestamp
}

// AddKey adds a new key/value pair to the fault's context if key is not empty.
//
// Parameters:
//   - fault: The fault to add the key/value pair to.
//   - key: The key to add.
//   - value: The value to add.
//
// Returns:
//   - bool: True if the key is empty or the fault is not nil, false otherwise.
func AddKey(fault flt.Fault, key string, value any) bool {
	if key == "" {
		return true
	}

	if fault == nil {
		return false
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	if base.context == nil {
		base.context = make(map[string]any)
	}

	base.context[key] = value

	return true
}

// GetValue gets the value of a key from the fault's context.
//
// Parameters:
//   - fault: The fault to get the value of the key from.
//   - key: The key to get the value of.
//
// Returns:
//   - any: The value of the key.
//   - flt.Fault: The fault that caused the error.
func GetValue(fault flt.Fault, key string) (any, flt.Fault) {
	if fault == nil {
		return nil, NewNilParameter("fault")
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	if len(base.context) == 0 {
		return nil, NewNoSuchKey(key)
	}

	value, ok := base.context[key]
	if !ok {
		return nil, NewNoSuchKey(key)
	}

	return value, nil
}

// ValueOf gets the value of a key from the fault's context and asserts that it is of type T.
//
// Parameters:
//   - fault: The fault to get the value of the key from.
//   - key: The key to get the value of.
//
// Returns:
//   - T: The value of the key.
//   - flt.Fault: The fault that caused the error.
func ValueOf[T any](fault flt.Fault, key string) (T, flt.Fault) {
	zero := *new(T)

	val, err := GetValue(fault, key)
	if err != nil {
		return zero, err
	}

	v, ok := val.(T)
	if !ok {
		err := NewNoSuchKey(key)

		SetSuggestions(err,
			fmt.Sprintf("key with type %T does not exist, but one of type %T was found", zero, val),
			"You may have forgotten to cast the value to the correct type or the desired key does not exist",
		)

		return zero, err
	}

	return v, nil
}

// SetValue sets the value of a key in the fault's context.
//
// Parameters:
//   - fault: The fault to set the value of the key in.
//   - key: The key to set the value of.
//   - value: The value to set.
//
// Returns:
//   - flt.Fault: The fault that caused the error.
func SetValue(fault flt.Fault, key string, value any) flt.Fault {
	if fault == nil {
		return NewNilParameter("fault")
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	if len(base.context) == 0 {
		return NewNoSuchKey(key)
	}

	_, ok = base.context[key]
	if !ok {
		return NewNoSuchKey(key)
	}

	base.context[key] = value

	return nil
}

// EditValue edits the value of a key in the fault's context.
//
// Parameters:
//   - fault: The fault to edit the value of the key in.
//   - key: The key to edit the value of.
//   - fn: The function to edit the value with.
//
// Returns:
//   - flt.Fault: The fault that caused the error.
func EditValue(fault flt.Fault, key string, fn func(v any) any) flt.Fault {
	if fault == nil {
		return NewNilParameter("fault")
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	v, err := GetValue(fault, key)
	if err != nil {
		return err
	}

	res := fn(v)

	err = SetValue(base, key, res)
	return err
}

// DeleteKey deletes a key from the fault's context.
//
// Parameters:
//   - fault: The fault to delete the key from.
//   - key: The key to delete.
//
// Returns:
//   - flt.Fault: The fault that caused the error.
func DeleteKey(fault flt.Fault, key string) {
	if fault == nil {
		return
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	delete(base.context, key)
}

// SetSuggestions sets the fault's suggestions; ignoring any empty suggestions.
//
// Parameters:
//   - fault: The fault to set the suggestions for.
//   - suggestions: The suggestions to set.
//
// Returns:
//   - bool: True if the suggestions were set, false otherwise. It only returns false
//     when there is at least one non-empty suggestion and the fault is nil.
func SetSuggestions(fault flt.Fault, suggestions ...string) bool {
	var count int

	for i := 0; i < len(suggestions); i++ {
		if suggestions[i] != "" {
			count++
		}
	}

	if count == 0 {
		return true
	}

	if fault == nil {
		return false
	}

	base, ok := Access[*flt.BaseFault](fault)
	if !ok {
		panic(flt.BadConstruction.Init())
	}

	filtered := make([]string, 0, count)

	for i := 0; i < len(suggestions); i++ {
		if suggestions[i] != "" {
			filtered = append(filtered, suggestions[i])
		}
	}

	base.suggestions = append(base.suggestions, filtered...)

	return true
}
