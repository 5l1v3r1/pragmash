package pragmash

import (
	"errors"
	"fmt"
	"os"
)

// StdInternal implements built-in commands that make the language usable.
type StdInternal struct {
	Variables map[string]Value
}

// Count returns the number of elements in a list.
func (s StdInternal) Count(args []string) Value {
	return NewNumberInt(int64(len(args)))
}

// Exit exits the current program with an optional exit code.
func (s StdInternal) Exit(args []Value) {
	if len(args) != 1 {
		os.Exit(0)
	} else {
		num, err := args[0].Number()
		if err != nil {
			os.Exit(1)
		}
		i := num.Int()
		if i != nil {
			os.Exit(int(i.Int64()))
		} else {
			os.Exit(int(num.Float()))
		}
	}
}

// Get gets a variable.
func (s StdInternal) Get(name string) (Value, error) {
	if val, ok := s.Variables[name]; ok {
		return val, nil
	}
	return nil, errors.New("variable undefined: " + name)
}

// Len returns the length of a string in bytes.
func (s StdInternal) Len(val string) Value {
	return NewNumberInt(int64(len(val)))
}

// Set sets a variable.
func (s StdInternal) Set(name string, val Value) {
	s.Variables[name] = val
}

// Throw throws an exception.
func (s StdInternal) Throw(args []Value) error {
	interfaceArgs := make([]interface{}, len(args))
	for i, x := range args {
		interfaceArgs[i] = x
	}
	return errors.New(fmt.Sprint(interfaceArgs...))
}
