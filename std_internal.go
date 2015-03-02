package pragmash

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// StdInternal implements built-in commands that make the language usable.
type StdInternal struct{}

// Call calls a function by expanding one or more lists of arguments.
func (_ StdInternal) Call(r Runner, name string,
	args ...[]*Value) (*Value, error) {
	totalCount := 0
	for _, x := range args {
		totalCount += len(x)
	}
	allArgs := make([]*Value, 0, totalCount)
	for _, x := range args {
		allArgs = append(allArgs, x...)
	}
	return r.RunCommand(name, allArgs)
}

// Count returns the number of elements in a list.
func (_ StdInternal) Count(args []*Value) *Value {
	count := int64(len(args))
	return NewValueNumber(NewNumberInt(count))
}

// Eval runs some pragmash code inside the current runner.
func (_ StdInternal) Eval(r Runner, code string) (*Value, error) {
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
	if val, bo := runnable.Run(r); bo == nil {
		return val, nil
	} else if bo.Type() == BreakoutTypeReturn {
		return bo.Value(), nil
	} else {
		return nil, errors.New(bo.Context() + ": " + bo.Error().Error())
	}
}

// Exec runs a pragmash script inside the current runner. It will be able to
// affect variables, throw exceptions, print to the console, etc.
func (_ StdInternal) Exec(r Runner, path string) (*Value, error) {
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
	if val, bo := runnable.Run(r); bo == nil {
		return val, nil
	} else if bo.Type() == BreakoutTypeReturn {
		return bo.Value(), nil
	} else {
		return nil, errors.New(bo.Context() + ": " + bo.Error().Error())
	}
}

// Exit exits the current program with an optional exit code.
func (_ StdInternal) Exit(args ...*Value) {
	if len(args) != 1 {
		os.Exit(0)
	} else {
		num, err := args[0].Number()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(int(num.Float()))
	}
}

// Len returns the length of a string in bytes.
func (_ StdInternal) Len(val string) int {
	return len(val)
}

// Pragmash runs a script with a given set of arguments in a new, standard
// runner. This is different from Exec because it isolates the variables of the
// new script and it sets its $DIR and $ARGV variables.
func (_ StdInternal) Pragmash(path string, args ...*Value) (*Value, error) {
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
	variables := CreateStandardVariables(path, args)
	runner := NewStdRunner(variables)

	// Run the file.
	if val, bo := runnable.Run(runner); bo == nil {
		return val, nil
	} else if bo.Type() == BreakoutTypeReturn {
		return bo.Value(), nil
	} else {
		return nil, errors.New(bo.Context() + ": " + bo.Error().Error())
	}
}

// Throw throws an exception.
func (_ StdInternal) Throw(args ...*Value) error {
	strArgs := make([]string, len(args))
	for i, x := range args {
		strArgs[i] = x.String()
	}
	return errors.New(strings.Join(strArgs, " "))
}
