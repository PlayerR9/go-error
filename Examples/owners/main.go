package main

import (
	"fmt"

	"github.com/PlayerR9/go-fault"
	"github.com/PlayerR9/go-fault/Examples/owners/internal"
)

var (
	// Owners is a map of owners and their pets.
	Owners map[string]string
)

func init() {
	Owners = make(map[string]string)

	Owners["Alice"] = "cat"
	Owners["Bob"] = "dog"
}

// PetOf returns the pet of a given owner.
//
// Parameters:
//   - name: The name of the owner.
//
// Returns:
//   - string: The pet of the owner.
//   - fault.Fault: The error if the owner was not found.
func PetOf(name string) (string, fault.Fault) {
	pet, ok := Owners[name]
	if !ok {
		return "", internal.NewErrKeyNotFound(name, "Owners")
	}

	return pet, nil
}

func main() {
	const (
		Mark string = "Mark"
	)

	pet, err := PetOf(Mark)
	if err == nil {
		_, err := fmt.Printf("%s's pet is a %s\n", Mark, pet)
		if err != nil {
			panic(err)
		}
	} else {
		lines := fault.LinesOf(err)

		for _, line := range lines {
			_, err := fmt.Println(line)
			if err != nil {
				panic(err)
			}
		}
	}
}
