package pragmash

import (
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// An Exception stores an error and context info (e.g. a line number) for that
// error.
type Exception struct {
	context string
	err     error
}

// NewException creates a new exception.
func NewException(context string, err error) Exception {
	return Exception{context, err}
}

// Array splits the error string by newlines and returns an array of exceptions,
// each with the same context but different lines from the original error.
func (e Exception) Array() []Value {
	str := e.String()
	if len(str) == 0 {
		return []Value{}
	}
	errorStrs := strings.Split(str, "\n")
	res := make([]Value, len(errorStrs))
	for i, x := range errorStrs {
		res[i] = NewException(e.context, errors.New(x))
	}
	return res
}

// Context returns the context string.
func (e Exception) Context() string {
	return e.context
}

// Error returns the error.
func (e Exception) Error() error {
	return e.err
}

// Number attempts to parse the error's string.
func (e Exception) Number() (Number, error) {
	return ParseNumber(e.String())
}

// String returns the error's string representation (i.e. e.Error().Error())
func (e Exception) String() string {
	return e.err.Error()
}

// A Number stores a numerical value.
type Number interface {
	Value

	// Float returns the float64 representation of the number. For integers with
	// large magnitudes, this may be +/- infinity.
	Float() float64

	// Int returns the big integer representation of the number if there is one.
	Int() *big.Int

	// Zero returns true if the number is zero.
	Zero() bool
}

type number struct {
	floating float64
	integer  *big.Int
}

func newNumberBig(b *big.Int) number {
	// Use big.Rat to convert the integer to a floating point.
	// NOTE: using float64(b.Int64()) will cause undefined behavior where there
	// need not be.
	rat := big.Rat{}
	rat.SetInt(b)
	f, _ := rat.Float64()

	return number{f, b}
}

func newNumberFloat(f float64) number {
	return number{f, nil}
}

func (n number) Array() []Value {
	return []Value{n}
}

func (n number) Context() string {
	return ""
}

func (n number) Float() float64 {
	return n.floating
}

func (n number) Int() *big.Int {
	return n.integer
}

func (n number) Number() (Number, error) {
	return n, nil
}

func (n number) String() string {
	if n.integer != nil {
		return n.integer.String()
	}
	return strconv.FormatFloat(n.floating, 'f', 10, 64)
}

func (n number) Zero() bool {
	if n.integer != nil {
		return n.integer.Cmp(big.NewInt(0)) == 0
	} else {
		return n.floating == 0
	}
}

// AddNumbers adds two numbers and returns the sum.
func AddNumbers(n1, n2 Number) Number {
	i1, i2 := n1.Int(), n2.Int()
	if i1 != nil && i2 != nil {
		return newNumberBig(big.NewInt(0).Add(i1, i2))
	} else {
		return newNumberFloat(n1.Float() + n2.Float())
	}
}

// DivideNumbers multiplies two numbers and returns the product.
// This returns an error if the second argument is zero.
func DivideNumbers(n1, n2 Number) (Number, error) {
	if n2.Zero() {
		return nil, errors.New("Division by zero.")
	}

	i1, i2 := n1.Int(), n2.Int()
	if i1 != nil && i2 != nil {
		rat := big.NewRat(0, 1)

		rat.SetFrac(i1, i2)
		if rat.IsInt() {
			// Special case where the division resulted in an integer.
			return newNumberBig(rat.Num()), nil
		}

		// Division resulted in a floating point.
		f, _ := rat.Float64()
		return newNumberFloat(f), nil
	} else {
		return newNumberFloat(n1.Float() / n2.Float()), nil
	}
}

// ExponentiateNumber raises a number to a given power.
func ExponentiateNumber(base, power number) Number {
	i1, i2 := base.Int(), power.Int()
	if i1 != nil && i2 != nil {
		return newNumberBig(big.NewInt(0).Exp(i1, i2, nil))
	} else {
		return newNumberFloat(math.Pow(base.Float(), power.Float()))
	}
}

// MultiplyNumbers multiplies two numbers and returns the product.
func MultiplyNumbers(n1, n2 Number) Number {
	i1, i2 := n1.Int(), n2.Int()
	if i1 != nil && i2 != nil {
		return newNumberBig(big.NewInt(0).Mul(i1, i2))
	} else {
		return newNumberFloat(n1.Float() + n2.Float())
	}
}

// ParseNumber parses a string and returns a number, or fails with an error.
func ParseNumber(s string) (Number, error) {
	// Parse it as a floating point.
	f, err := strconv.ParseFloat(s, 64)
	if strings.Contains(s, ".") {
		if err != nil {
			return nil, err
		}
		return number{f, nil}, nil
	}

	// Parse it as a big int.

	// NOTE: if the number was HUGE, ParseFloat() would have returned an error
	// even though our big.Int will be fine. Thus, we let the error slide.

	num := big.NewInt(0)
	if _, ok := num.SetString(s, 10); !ok {
		return nil, errors.New("Invalid integer: " + s)
	}
	return number{f, num}, nil
}

// SubtractNumbers subtracts two numbers and returns the difference.
func SubtractNumbers(n1, n2 Number) Number {
	i1, i2 := n1.Int(), n2.Int()
	if i1 != nil && i2 != nil {
		return newNumberBig(big.NewInt(0).Sub(i1, i2))
	} else {
		return newNumberFloat(n1.Float() - n2.Float())
	}
}

// A Runnable is a generic interface which can execute on a given Runner and
// return a value or runtime exception.
type Runnable interface {
	Run(r Runner) (string, Exception)
}

// A Runner is a generic interface which can run a commands.
type Runner interface {
	RunCommand(name string, args []Value) (Value, error)
}

// A Value is a read-only variable value.
type Value interface {
	// Array returns the array representation of the value.
	Array() []Value

	// Context returns the context of the value. This is useful if the value is
	// an exception. In most cases, this should be an empty string.
	Context() string

	// Number returns the numerical representation of the value, or an error if
	// the value is not a number.
	Number() (Number, error)

	// String returns the textual representation of the value.
	String() string
}
