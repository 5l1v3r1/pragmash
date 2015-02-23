package pragmash

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// StdInternal implements built-in commands that make the language usable.
type StdInternal struct {
	Variables map[string]Value
	Runner    *Runner
}

// NewStdInternal creates a StdInternal with some default variables.
func NewStdInternal() StdInternal {
	return StdInternal{map[string]Value{}, nil}
}

// Count returns the number of elements in a list.
func (_ StdInternal) Count(args []string) Value {
	return NewNumberInt(int64(len(args)))
}

// Eval runs some pragmash code inside the current runner.
func (s StdInternal) Eval(code string) (Value, error) {
	if s.Runner == nil || *s.Runner == nil {
		return nil, errors.New("no Runner")
	}
	
	lines, contexts, err := TokenizeString(code)
	if err != nil {
		return nil, err
	}
	
	// Update the contexts to reflect that we're in an eval.
	for i, x := range contexts {
		contexts[i] = x + " in eval"
	}
	
	// Generate the runner and run it
	runnable, err := ScanAll(lines, contexts)
	if err != nil {
		return nil, err
	}
	if val, err := runnable.Run(*s.Runner); err != nil {
		return nil, errors.New(err.Context() + ": " + err.String())
	} else {
		return val, nil
	}
}

// Exec runs a pragmash script inside the current runner. It will be able to
// affect variables, throw exceptions, print to the console, etc.
func (s StdInternal) Exec(path string) (Value, error) {
	if s.Runner == nil || *s.Runner == nil {
		return nil, errors.New("no Runner")
	}
	
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines, contexts, err := TokenizeString(string(contents))
	if err != nil {
		return nil, err
	}
	
	// Update the contexts to include the path
	for i, x := range contexts {
		contexts[i] = x + " in " + path
	}
	
	// Generate the runner and run it
	runnable, err := ScanAll(lines, contexts)
	if err != nil {
		return nil, err
	}
	if val, err := runnable.Run(*s.Runner); err != nil {
		return nil, errors.New(err.Context() + ": " + err.String())
	} else {
		return val, nil
	}
}

// Exit exits the current program with an optional exit code.
func (_ StdInternal) Exit(args []Value) {
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
func (_ StdInternal) Len(val string) Value {
	return NewNumberInt(int64(len(val)))
}

// Pragmash runs a script with a given set of arguments in a new, standard
// runner. This is different from Exec because it isolates the variables of the
// new script and it sets its $DIR and $ARGV variables.
func (_ StdInternal) Pragmash(args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, errors.New("missing file path")
	}
	
	path := args[0].String()
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines, contexts, err := TokenizeString(string(contents))
	if err != nil {
		return nil, err
	}
	
	// Update the contexts to include the path
	for i, x := range contexts {
		contexts[i] = x + " in " + path
	}
	
	// Generate the runnable
	runnable, err := ScanAll(lines, contexts)
	if err != nil {
		return nil, err
	}
	
	// Generate the runner
	strArgs := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		strArgs[i-1] = args[i].String()
	}
	variables := map[string]Value{
		"DIR": NewHybridValueString(filepath.Dir(path)),
		"ARGV": NewHybridValueString(strings.Join(strArgs, "\n")),
	}
	runner := NewStdRunner(variables)
	
	// Run the file.
	if val, err := runnable.Run(runner); err != nil {
		return nil, errors.New(err.Context() + ": " + err.String())
	} else {
		return val, nil
	}
}

// Set sets a variable.
func (s StdInternal) Set(name string, val Value) {
	s.Variables[name] = val
}

// Throw throws an exception.
func (_ StdInternal) Throw(args []Value) error {
	strArgs := make([]string, len(args))
	for i, x := range args {
		strArgs[i] = x.String()
	}
	return errors.New(strings.Join(strArgs, " "))
}
