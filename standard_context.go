package pragmash

import (
	"errors"
	"fmt"
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
		"echo": res.Echo,
		"get":  res.Get,
		"puts": res.Puts,
		"set":  res.Set,
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
