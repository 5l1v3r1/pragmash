package pragmash

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// A CommandFunc handles a specific command.
type CommandFunc func([]string) (string, error)

// A StandardContext is an extensible context with built-in functionality.
type StandardContext struct {
	Commands  map[string]CommandFunc
	Variables map[string]string
}

// NewStandardContext creates a new standard context and returns it.
func NewStandardContext() *StandardContext {
	res := &StandardContext{Variables: map[string]string{}}
	res.Commands = map[string]CommandFunc{
		"echo":  res.Echo,
		"get":   res.Get,
		"puts":  res.Puts,
		"range": res.Range,
		"set":   res.Set,
		"throw": res.Throw,
	}
	return res
}

// Echo returns a space-delimited version of the arguments.
func (s *StandardContext) Echo(args []string) (string, error) {
	return strings.Join(args, " "), nil
}

// Get returns the value of a given variable or an error if the variable is
// undefined.
func (s *StandardContext) Get(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("Missing arguments to 'get' command.")
	}
	if x, ok := s.Variables[args[0]]; ok {
		return x, nil
	} else {
		return "", errors.New("Undefined variable: " + args[1])
	}
}

// Puts prints text to the console and returns an empty string.
func (s *StandardContext) Puts(args []string) (string, error) {
	fmt.Println(strings.Join(args, " "))
	return "", nil
}

// Range returns a newline-delimited list of integers in a given range.
func (s *StandardContext) Range(args []string) (string, error) {
	// Validate argument count.
	if len(args) == 0 || len(args) > 3 {
		return "", errors.New("Range takes 1, 2, or 3 arguments, got " +
			strconv.Itoa(len(args)))
	}
	
	// Parse arguments.
	numArgs := make([]int, len(args))
	for i, x := range args {
		var err error
		numArgs[i], err = strconv.Atoi(x)
		if err != nil {
			return "", err
		}
	}
	start := 0
	end := numArgs[0]
	step := 1
	if len(args) >= 2 {
		start, end = end, numArgs[1]
	}
	if len(args) == 3 {
		step = numArgs[2]
		if step == 0 {
			return "", errors.New("Step cannot be zero.")
		}
	}
	
	// Generate the range.
	if step > 0 {
		if end < start {
			return "", nil
		}
		res := ""
		for i := start; i < end; i += step {
			res += strconv.Itoa(i)
			if i + step < end {
				res += "\n"
			}
		}
		return res, nil
	} else {
		if end > start {
			return "", nil
		}
		res := ""
		for i := start; i > end; i += step {
			res += strconv.Itoa(i)
			if i + step > end {
				res += "\n"
			}
		}
		return res, nil
	}
}

// Run runs a command to satisfy the Context interface.
func (s *StandardContext) Run(command string, args []string) (string, error) {
	if cmd, ok := s.Commands[command]; ok {
		return cmd(args)
	} else {
		return "", errors.New("Unknown command: " + command)
	}
}

// Set sets a variable's value.
func (s *StandardContext) Set(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Missing arguments to 'set' command.")
	}
	s.Variables[args[0]] = args[1]
	return "", nil
}

// Throw generates an error.
func (s *StandardContext) Throw(args []string) (string, error) {
	return "", errors.New(strings.Join(args, " "))
}
