package pragmash

import (
	"path/filepath"
)

// StdAll implements the methods corresponding to the standard library.
type StdAll struct {
	StdArray
	StdFs
	StdInternal
	StdIo
	StdMath
	StdOps
	StdString
	StdTime
}

// CreateStandardVariables generates the set of standard variables for a given
// script and set of arguments.
func CreateStandardVariables(script string, args []*Value) map[string]*Value {
	cleanPath, err := filepath.Abs(filepath.Clean(script))
	if err != nil {
		cleanPath = filepath.Clean(script)
	}
	return map[string]*Value{
		"ARGV":    NewValueArray(args),
		"DIR":     NewValueString(filepath.Dir(cleanPath)),
		"SCRIPT":  NewValueString(cleanPath),
		"VERSION": NewValueString(Version()),
	}
}

// NewStdRunner returns a Runner which implements the standard library.
func NewStdRunner(variables map[string]*Value) Runner {
	runner := NewReflectRunner(StdAll{}, OperatorRewrites)

	// Copy variables if necessary.
	if variables != nil {
		for name, value := range variables {
			runner.variables[name] = value
		}
	}

	return runner
}
