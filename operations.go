package fault

import "fmt"

// AddKey adds a new key/value pair to the fault's context if key is not empty.
//
// Parameters:
//   - fault: The fault to add the key/value pair to.
//   - key: The key to add.
//   - value: The value to add.
//
// Returns:
//   - bool: True if the key is empty or the fault is not nil, false otherwise.
func AddKey(fault Fault, key string, value any) bool {
	if key == "" {
		return true
	}

	if fault == nil {
		return false
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
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
//   - Fault: The fault that caused the error.
func GetValue(fault Fault, key string) (any, Fault) {
	if fault == nil {
		return nil, NewNilParameter("fault")
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
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
//   - Fault: The fault that caused the error.
func ValueOf[T any](fault Fault, key string) (T, Fault) {
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
//   - Fault: The fault that caused the error.
func SetValue(fault Fault, key string, value any) Fault {
	if fault == nil {
		return NewNilParameter("fault")
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
	}

	if len(base.context) == 0 {
		return NewNoSuchKey(key)
	}

	_, ok := base.context[key]
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
//   - Fault: The fault that caused the error.
func EditValue(fault Fault, key string, fn func(v any) any) Fault {
	if fault == nil {
		return NewNilParameter("fault")
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
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
//   - Fault: The fault that caused the error.
func DeleteKey(fault Fault, key string) {
	if fault == nil {
		return
	}

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
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
func SetSuggestions(fault Fault, suggestions ...string) bool {
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

	base := get_base(fault)
	if base == nil {
		panic(BadConstruction.Init())
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
