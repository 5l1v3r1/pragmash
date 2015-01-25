package pragmash

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
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
		"get":   res.Get,
		"gets":  res.Gets,
		"*":     res.Multiply,
		"print": res.Print,
		"puts":  res.Puts,
		"range": res.Range,
		"set":   res.Set,
		"-":     res.Subtract,
		"throw": res.Throw,
	}
	return res
}

// Add adds a bunch of big integers or floating points.
func (s *StandardContext) Add(args []string) (string, error) {
	if len(args) == 0 {
		return "0", nil
	}
	if numsUseFloat(args) {
		// Use floating point.
		floats, err := numsParseFloats(args)
		if err != nil {
			return "", err
		}
		sum := 0.0
		for _, x := range floats {
			sum += x
		}
		return strconv.FormatFloat(sum, 'f', 10, 64), nil
	} else {
		// Use big integer.
		ints, err := numsParseInts(args)
		if err != nil {
			return "", err
		}
		for i := 1; i < len(ints); i++ {
			ints[0].Add(ints[0], ints[i])
		}
		return ints[0].String(), nil
	}
}

// Divide divides one floating point by another one.
func (s *StandardContext) Divide(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Division takes exactly two arguments.")
	}
	
	floats, err := numsParseFloats(args)
	if err != nil {
		return "", err
	}
	if floats[1] == 0.0 {
		return "", errors.New("Division by zero.")
	}
	return strconv.FormatFloat(floats[0] / floats[1], 'f', 10, 64), nil
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

// Gets reads a line from the console and returns it without a newline
// character.
func (s *StandardContext) Gets(args []string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), nil
}

// Multiply multiplies a list of big integers or floating points.
func (s *StandardContext) Multiply(args []string) (string, error) {
	if len(args) == 0 {
		return "1", nil
	}
	if numsUseFloat(args) {
		// Use floating point.
		floats, err := numsParseFloats(args)
		if err != nil {
			return "", err
		}
		product := 1.0
		for _, x := range floats {
			product *= x
		}
		return strconv.FormatFloat(product, 'f', 10, 64), nil
	} else {
		// Use big integer.
		ints, err := numsParseInts(args)
		if err != nil {
			return "", err
		}
		for i := 1; i < len(ints); i++ {
			ints[0].Mul(ints[0], ints[i])
		}
		return ints[0].String(), nil
	}
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

// Subtract performs subtraction of big integers or floating points.
func (s *StandardContext) Subtract(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("Subtraction takes exactly two arguments.")
	}
	if numsUseFloat(args) {
		// Use floating point.
		floats, err := numsParseFloats(args)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(floats[0] - floats[1], 'f', 10, 64), nil
	} else {
		// Use big integer.
		ints, err := numsParseInts(args)
		if err != nil {
			return "", err
		}
		ints[0].Sub(ints[0], ints[1])
		return ints[0].String(), nil
	}
}

// Throw generates an error.
func (s *StandardContext) Throw(args []string) (string, error) {
	return "", errors.New(strings.Join(args, " "))
}

func numsUseFloat(nums []string) bool {
	for _, x := range nums {
		if strings.Contains(x, ".") {
			return true
		}
	}
	return false
}

func numsParseFloats(nums []string) ([]float64, error) {
	res := make([]float64, len(nums))
	for i, x := range nums {
		var err error
		res[i], err = strconv.ParseFloat(x, 64)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func numsParseInts(nums []string) ([]*big.Int, error) {
	res := make([]*big.Int, len(nums))
	for i, x := range nums {
		num := big.NewInt(0)
		if _, ok := num.SetString(x, 10); !ok {
			return nil, errors.New("Invalid integer: " + x)
		}
		res[i] = num
	}
	return res, nil
}
