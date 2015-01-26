package pragmash

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
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
		"=":     res.Equal,
		"exit":  res.Exit,
		"get":   res.Get,
		"[]":    res.GetAt,
		"gets":  res.Gets,
		">=":    res.GreaterEqual,
		">":     res.GreaterThan,
		"len":   res.Len,
		"<=":    res.LessEqual,
		"<":     res.LessThan,
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

// Equal returns "true" if all its arguments are equal, or "false" otherwise.
func (s *StandardContext) Equal(args []string) (string, error) {
	for i := 1; i < len(args); i++ {
		if args[i] != args[0] {
			return "", nil
		}
	}
	return "true", nil
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
		return "", errors.New("Undefined variable: " + args[1])
	}
}

// GetAt returns an entry in an array, or throws an error if the index is out of
// bounds.
func (s *StandardContext) GetAt(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("The subscirpt operator takes two arguments.")
	}
	if len(args[0]) == 0 {
		return "", errors.New("List is empty.")
	}
	list := strings.Split(args[0], "\n")
	idx, err := strconv.Atoi(args[1])
	if err != nil {
		return "", err
	}
	if idx < 0 || idx >= len(list) {
		return "", errors.New("Index out of bounds: " + args[1])
	}
	return list[idx], nil
}

// Gets reads a line from the console and returns it without a newline
// character.
func (s *StandardContext) Gets(args []string) (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), nil
}

// GreaterEqual compares two integers or floating points and returns true if the
// first argument is greater than or equal to the second.
func (s *StandardContext) GreaterEqual(args []string) (string, error) {
	cmp, err := compareNums(args)
	if err != nil {
		return "", err
	}
	if cmp >= 0 {
		return "true", nil
	} else {
		return "", nil
	}
}

// GreaterThan compares two integers or floating points and returns true if the
// first argument is greater than the second.
func (s *StandardContext) GreaterThan(args []string) (string, error) {
	cmp, err := compareNums(args)
	if err != nil {
		return "", err
	}
	if cmp > 0 {
		return "true", nil
	} else {
		return "", nil
	}
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

// LessEqual compares two integers or floating points and returns true if the
// first argument is less than or equal to the second.
func (s *StandardContext) LessEqual(args []string) (string, error) {
	cmp, err := compareNums(args)
	if err != nil {
		return "", err
	}
	if cmp <= 0 {
		return "true", nil
	} else {
		return "", nil
	}
}

// LessThan compares two integers or floating points and returns true if the
// first argument is less than the second.
func (s *StandardContext) LessThan(args []string) (string, error) {
	cmp, err := compareNums(args)
	if err != nil {
		return "", err
	}
	if cmp < 0 {
		return "true", nil
	} else {
		return "", nil
	}
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

// Not inverts a conditional.
func (s *StandardContext) Not(args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	if len(args) == 1 {
		if len(args[0]) == 0 {
			return "true", nil
		} else {
			return "", nil
		}
	}
	for i := 1; i < len(args); i++ {
		if args[i] != args[0] {
			return "true", nil
		}
	}
	return "", nil
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

// Write writes a string to a file.
func (s *StandardContext) Write(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.New("The write command expects two arguments.")
	}
	path := args[0]
	data := args[1]
	if err := ioutil.WriteFile(path, []byte(data), os.FileMode(0600));
		err != nil {
		return "", err
	}
	return "", nil
}

func compareNums(args []string) (int, error) {
	if len(args) != 0 {
		return 0, errors.New("Comparisons take two arguments, but got " +
			strconv.Itoa(len(args)))
	}
	if numsUseFloat(args) {
		// Use floating point.
		floats, err := numsParseFloats(args)
		if err != nil {
			return 0, err
		}
		if floats[0] < floats[1] {
			return -1, nil
		} else if floats[0] == floats[1] {
			return 0, nil
		}
		return 1, nil
	} else {
		// Use big integer.
		ints, err := numsParseInts(args)
		if err != nil {
			return 0, err
		}
		return ints[0].Cmp(ints[1]), nil
	}
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

func numsUseFloat(nums []string) bool {
	for _, x := range nums {
		if strings.Contains(x, ".") {
			return true
		}
	}
	return false
}
