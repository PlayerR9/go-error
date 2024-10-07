# go-fault
A go package that implements a more "complex" way to handle errors in Go. (Alternative to the standard library)


## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Overview](#overview)
3. [Installation](#installation)
4. [Usage](#usage)


## Overview

The `go-fault` package aims to "bring" the error handling capabilities of the SD programming language to Go; whilst maintaining interopability with the existing standard library.


***What Is a Fault?***

In the SD programming language, errors are handled by the fault system. Simply put, a **fault** is a data type that carries information about an error in a standardized way. This means that, not only it allows for more clear and descriptive error messages, but it also allows for generalization and extensibility.

A fault must carry at least these information:
- **Level:** The level (or severity) of the fault allows the caller to know the "importance" of the fault.
- **Code:** The code of the fault allows the caller to identify the broader category of the fault such as "BadParameter", "OperationFail", etc.
- **Message:** The message of the fault allows to narrow down the fault to its root cause; as well as to inform the caller about the nature of the fault.
- **Timestamp:** The timestamp of the fault specifies when the fault occurred.
- **Suggestions:** The suggestions of the fault describes one or more possible solutions or actions that can be taken to resolve the fault.


When displayed, a fault is formatted as:
```console
[<level>] (<code>) <msg>

Occurred at: <timestamp>
// Additional information...
```
Where;
- `<level>`: The level of the fault.
- `<code>`: The code of the fault.
- `<msg>`: The message of the fault.
- `<timestamp>`: The timestamp of the fault.
- Additional information are the additional information of the fault.


## Installation

To use the `go-fault` package, simply import it in your project as it doesn't need any additional dependencies. In order to get the package in your project, use the following command:
```bash
$ go get github.com/PlayerR9/go-fault
```

And, once it's installed, you can use it in your project with the import statement:
```go
import "github.com/PlayerR9/go-fault"
```


## Usage

***How to Define a Fault?***

Generally speaking, a custom fault is declared as follows:
```go
import "github.com/PlayerR9/go-fault"

type MyFault struct {
   fault.Fault

   // Additional fields...
}

func (e MyFault) Embeds() fault.Fault {
   return e.Fault
}

func (e MyFault) InfoLines() []string {
   var lines []string

   // Write here any additional information of the fault.
   // IMPORTANT: No not call e.Fault.InfoLines()! 

   return lines
}
```

Here, `MyFault` is a fault type that embeds another fault (often referred to as the "base") and implements the `fault.Fault` interface.

However, as you may have already noticed, the fault itself does not implement the `Error() string` method. This is due to the fact that it is up to the fault's base to implement it.


***How to Create a Fault?***

Any constructor for a custom fault looks like the following:
```go
func NewMyFault(/*Arguments go here*/) *MyFault {
   base := // Call the constructor of the fault's base

   return &MyFault{
      Fault: base,

      // Additional fields...
   }
}
```


## Example

Here's an example:
```go
package internal

import (
	"strconv"

	"github.com/PlayerR9/go-fault"
	"github.com/PlayerR9/go-fault/faults"
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
	base := fault.New(faults.OperationFail, "the specified key was not found")

	return &ErrKeyNotFound{
		Fault:   base,
		Key:     key,
		SetName: set_name,
	}
}
```


```go
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
```

```bash
$ go run main.go

# [ERROR] (OperationFail) the specified key was not found.

# - Occurred at: 2024-10-07 09:47:12.402688435 +0200 CEST m=+0.000018469
# - Key: "Mark"
# - Set name: "Owners"
```