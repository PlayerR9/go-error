# go-error
A go package that implements a more "complex" way to handle errors in Go. (Alternative to the standard library)


## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Overview](#overview)
3. [Installation](#installation)
4. [Usage](#usage)


## Overview

The `go-error` package aims to "bring" the error handling capabilities of the SD programming language to Go; whilst maintaining interopability with the existing standard library.


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

To use the `go-error` package, simply import it in your project as it doesn't need any additional dependencies.


## Usage

***How to Define a Fault?***

Generally speaking, a custom fault is declared as follows:
```go
import "github.com/PlayerR9/go-error"

type MyFault struct {
   errs.Fault

   // Additional fields...
}

func (e MyFault) Embeds() errs.Fault {
   return e.Fault
}

func (e MyFault) WriteInfo(w errs.Writer) (int, errs.Fault) {
   var total int

   // Write here any additional information of the fault.
   // IMPORTANT: No not call e.Fault.WriteInfo()! 

   return total, nil 
}
```

Here, `MyFault` is a fault type that embeds another fault (often referred to as the "base") and implements the `errs.Fault` interface.

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
import (
   "strconv"

   "github.com/PlayerR9/go-error"
)

type ErrKeyNotFound struct {
   errs.Fault

   Key      string
   SetName  string
}

func (e ErrKeyNotFound) Embeds() errs.Fault {
   return e.Fault
}

func (e ErrKeyNotFound) WriteInfo(w errs.Writer) (int, errs.Fault) {
   var total int

   data := []byte("- Key: " + e.Key + "\n- Set name: " + strconv.Quote(e.SetName))

   n, err := errs.Write(w, data)
   total += n

   return total, err
}

func NewErrKeyNotFound(key string, set_name string) *ErrKeyNotFound {
   base := errs.New(errs.OperationFail, "the specified key was not found")

   return &ErrKeyNotFound{
      Fault:      base,
      Key:        key,
      SetName:    set_name,
   }
}
```


```go
var (
   Owners map[string]string
)

func init() {
   Owners = make(map[string]string)
   Owners["Alice"] = "cat"
   Owners["Bob"] = "dog"
}

func PetOf(name string) (string, errs.Fault) {
   pet, ok := Owners[name]
   if !ok {
      return "", NewErrKeyNotFound(name, "Owners")
   }

   return pet, nil
}

func main() {
   pet, err := PetOf("Mark")
   if err == nil {
      _, err := fmt.Println("Mark's pet is a", pet)
      if err != nil {
         panic(err)
      }

      os.Exit(0)
   }
   
   var builder strings.Builder

   w := errs.WriterOf(&builder)

   _, err = errs.WriteFault(w, err)
   if err != nil {
      panic(err)
   }

   _, err = fmt.Println(builder.String())
   if err != nil {
      panic(err)
   }
}

```

```bash
$ go run main.go

# [ERROR] (OperationFail) the specified key was not found
# 
# 
# Occurred at: 2022-11-10T15:20:00.000000000+00:00
# Key: Mark
# Set name: "Owners"
```