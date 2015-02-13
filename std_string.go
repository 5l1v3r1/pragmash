package pragmash

import (
	"fmt"
)

// StdString implements ways of manipulating or creating strings
type StdString struct{}

// Echo joins its arguments with spaces.
func (s StdString) Echo(args []Value) Value {
	interfaceArgs := make([]interface{}, len(args))
	for i, x := range args {
		interfaceArgs[i] = x
	}
	return StringValue(fmt.Sprint(interfaceArgs...))
}
