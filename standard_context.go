package pragmash

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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
		"+":     res.Add,
		"/":     res.Divide,
		"echo":  res.Echo,
		"=":     res.Equal,
		"exit":  res.Exit,
		"get":   res.Get,
		"[]":    res.GetAt,
		"gets":  res.Gets,
		">=":    res.GreaterEqual,
		">":     res.GreaterThan,
		"join":  res.Join,
		"len":   res.Len,
		"<=":    res.LessEqual,
		"<":     res.LessThan,
		"match": res.Match,
		"*":     res.Multiply,
		"!":     res.Not,
		"print": res.Print,
		"puts":  res.Puts,
		"range": res.Range,
		"read":  res.Read,
		"set":   res.Set,
		"-":     res.Subtract,
		"throw": res.Throw,
		"write": res.Write,
	}
	return res
}

// Echo returns a space-delimited version of the arguments.
func (s *StandardContext) Echo(args []string) (string, error) {
	return strings.Join(args, " "), nil
}

// Exit exits the program with an optional return code.
func (s *StandardContext) Exit(args []string) (string, error) {
	if len(args) == 0 {
		os.Exit(0)
	} else if len(args) == 1 {
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return "", err
		}
		os.Exit(num)
	}
	return "", errors.New("Exit command takes 0 or 1 argument(s).")
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
		return "", errors.New("Undefined variable: " + args[0])
	}
}

// Gets reads a line from the console and returns it without a newline
// character.
func (s *StandardContext) Gets(args []string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), nil
}

// Join joins strings by appending them.
func (s *StandardContext) Join(args []string) (string, error) {
	res := ""
	for _, x := range args {
		res += x
	}
	return res, nil
}

// Len returns the number of lines in a string, or 0 if it's empty.
func (s *StandardContext) Len(args []string) (string, error) {
	count := 0
	for _, arg := range args {
		if len(arg) == 0 {
			continue
		}
		count += strings.Count(arg, "\n") + 1
	}
	return strconv.Itoa(count), nil
}

// Match runs a regular expression on a string.
func (s *StandardContext) Match(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("The match command takes two arguments.")
	}
	r, err := regexp.Compile(args[0])
	if err != nil {
		return "", err
	}
	res := r.FindAllStringSubmatch(args[1], -1)
	resStr := ""
	for i, x := range res {
		for j, y := range x {
			if i != 0 || j != 0 {
				resStr += "\n" + y
			}
		}
	}
	return resStr, nil
}

// Print prints text to the console without a newline and returns an empty
// string.
func (s *StandardContext) Print(args []string) (string, error) {
	fmt.Print(strings.Join(args, " "))
	return "", nil
}

// Puts prints text to the console with a newline and returns an empty string.
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
			if i+step < end {
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
			if i+step > end {
				res += "\n"
			}
		}
		return res, nil
	}
}

// Read reads the contents of a file or URL.
func (s *StandardContext) Read(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("The read command expects one argument.")
	}

	// Read a web URL if applicable.
	if strings.HasPrefix(args[0], "http://") ||
		strings.HasPrefix(args[0], "https://") {
		resp, err := http.Get(args[0])
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	}

	// Read a path.
	contents, err := ioutil.ReadFile(args[0])
	if err != nil {
		return "", err
	}
	return string(contents), nil
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

// Write writes a string to a file.
func (s *StandardContext) Write(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("The write command expects two arguments.")
	}
	path := args[0]
	data := args[1]
	if err := ioutil.WriteFile(path, []byte(data), os.FileMode(0600)); err != nil {
		return "", err
	}
	return "", nil
}
