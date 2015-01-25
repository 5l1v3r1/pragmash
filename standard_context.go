package pragmash

import (
	"errors"
	"fmt"
	"strings"
)

type CommandFunc func([]string) (string, error)

type StandardContext struct {
	Commands  map[string]CommandFunc
	Variables map[string]string
}

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

func (s *StandardContext) Echo(args []string) (string, error) {
	return strings.Join(args, " "), nil
}

func (s *StandardContext) Get(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("Missing arguments to 'get' command.")
	}
	if x, ok := s.Variables[args[1]]; ok {
		return x, nil
	} else {
		return "", nil
	}
}

func (s *StandardContext) Puts(args []string) (string, error) {
	fmt.Println(strings.Join(args, " "))
	return "", nil
}

func (s *StandardContext) Set(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Missing arguments to 'set' command.")
	}
	s.Variables[args[1]] = args[2]
	return "", nil
}
