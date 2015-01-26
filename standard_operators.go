package pragmash

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

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
	return strconv.FormatFloat(floats[0]/floats[1], 'f', 10, 64), nil
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
		return strconv.FormatFloat(floats[0]-floats[1], 'f', 10, 64), nil
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

func compareNums(args []string) (int, error) {
	if len(args) != 2 {
		return 0, errors.New("Comparisons take 2 arguments, but got " +
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
