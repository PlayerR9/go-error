package internal

import (
	"strconv"

	"github.com/PlayerR9/go-fault"
)

// ErrKeyNotFound is an error that indicates that the specified key was not found.
type ErrKeyNotFound struct {
	fault.Fault

	// Key is the key that was not found.
	Key string

	// SetName is the name of the set that is being accessed.
	SetName string
}

// Embeds implements the Fault interface.
func (e ErrKeyNotFound) Embeds() fault.Fault {
	return e.Fault
}

// InfoLines implements the Fault interface.
//
// Format:
//
//	"- Key: <key>"
//	"- Set name: <set_name>"
//
// Where:
//   - <key>: The key that was not found.
//   - <set_name>: The name of the set that is being accessed.
func (e ErrKeyNotFound) InfoLines() []string {
	lines := make([]string, 0, 2)

	lines = append(lines, "- Key: "+strconv.Quote(e.Key))
	lines = append(lines, "- Set name: "+strconv.Quote(e.SetName))

	return lines
}

// NewErrKeyNotFound returns an error that indicates that the specified key was not found.
//
// Parameters:
//   - key: The key that was not found.
//   - set_name: The name of the set that is being accessed.
//
// Returns:
//   - *ErrKeyNotFound: The error that indicates that the specified key was not found.
//     Never returns nil.
func NewErrKeyNotFound(key string, set_name string) *ErrKeyNotFound {
	base := fault.New(fault.OperationFailed, "the specified key was not found")

	return &ErrKeyNotFound{
		Fault:   base,
		Key:     key,
		SetName: set_name,
	}
}
