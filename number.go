package pragmash

import (
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Number struct {
	isInteger bool
	floating  float64
	integer   big.Int
}

// NewNumberBig returns an object which implements Number and represents the
// given integer.
func NewNumberBig(b *big.Int) *Number {
	// Use big.Rat to convert the integer to a floating point.
	// NOTE: using float64(b.Int64()) will not work since floating points can
	// be larger than 2^64.
	rat := big.Rat{}
	rat.SetInt(b)
	f, _ := rat.Float64()
	res := Number{true, f}
	res.integer.Set(b)
	return &res
}

// NewNumberFloat generates a Number which probably will not have an integer
// representation.
func NewNumberFloat(f float64) *Number {
	rat := big.Rat{}
	rat.SetFloat64(f)
	if rat.IsInt() {
		return NewNumberInt(f.Num())
	}
	return &Number{false, f}
}

// NewNumberInt generates a Number with an integer.
func NewNumberInt(i int64) *Number {
	return NewNumberBig(big.NewInt(i))
}

// Float returns the floating point representation of the number.
func (n *Number) Float() float64 {
	return n.floating
}

// Int returns the integer representation of the number or nil.
// You should not modify this value.
func (n *Number) Int() *big.Int {
	if !n.isInteger {
		return nil
	}
	return &n.integer
}

// IsInt returns true if the number is an integer.
func (n *Number) IsInt() bool {
	return n.isInteger
}

// Zero returns true if the number is zero.
func (n *Number) Zero() bool {
	if n.integer != nil {
		return n.integer.Cmp(&big.Int{}) == 0
	} else {
		return n.floating == 0
	}
}

// AddNumbers adds two numbers and returns the sum.
func AddNumbers(n1, n2 *Number) *Number {
	i1, i2 := n1.Int(), n2.Int()
	if i1 != nil && i2 != nil {
		return NewNumberBig(big.NewInt(0).Add(i1, i2))
	} else {
		return NewNumberFloat(n1.Float() + n2.Float())
	}
}

// CompareNumbers returns -1 if n1 < n2, 0 if n1 == n2, or 1 if n1 > n2.
func CompareNumbers(n1, n2 *Number) int {
	i1, i2 := n1.Int(), n2.Int()

	// Check if we need to use floating points.
	if i1 == nil || i2 == nil {
		f1, f2 := n1.Float(), n2.Float()
		if f1 < f2 {
			return -1
		} else if f1 > f2 {
			return 1
		} else {
			return 0
		}
	}

	return i1.Cmp(i2)
}

// DivideNumbers multiplies two numbers and returns the product.
// This returns an error if the second argument is zero.
func DivideNumbers(n1, n2 *Number) (*Number, error) {
	if n2.Zero() {
		return nil, errors.New("division by zero")
	}

	i1, i2 := n1.Int(), n2.Int()

	// See if we need to return a floating point.
	if i1 == nil || i2 == nil {
		return NewNumberFloat(n1.Float() / n2.Float()), nil
	}

	// Use a big rational number to see if we can return an integer.
	rat := &big.Rat{}
	rat.SetFrac(i1, i2)

	if rat.IsInt() {
		// Special case where the division resulted in an integer.
		return NewNumberBig(rat.Num()), nil
	}

	// Division resulted in a floating point.
	f, _ := rat.Float64()
	return NewNumberFloat(f), nil
}

// ExponentiateNumber raises a number to a given power.
func ExponentiateNumber(base, power *Number) *Number {
	i1, i2 := base.Int(), power.Int()
	if i1 != nil && i2 != nil {
		return NewNumberBig(big.NewInt(0).Exp(i1, i2, nil))
	} else {
		return NewNumberFloat(math.Pow(base.Float(), power.Float()))
	}
}

// MultiplyNumbers multiplies two numbers and returns the product.
func MultiplyNumbers(n1, n2 *Number) *Number {
	i1, i2 := n1.Int(), n2.Int()
	if i1 != nil && i2 != nil {
		return NewNumberBig(big.NewInt(0).Mul(i1, i2))
	} else {
		return NewNumberFloat(n1.Float() * n2.Float())
	}
}

// ParseNumber parses a string and returns a number, or fails with an error.
func ParseNumber(s string) (*Number, error) {
	// Parse it as a floating point.
	f, err := strconv.ParseFloat(s, 64)
	if strings.Contains(s, ".") {
		if err != nil {
			return nil, err
		}
		return NewNumberFloat(f)
	}

	// Parse it as a big int.

	// NOTE: if the number was HUGE, ParseFloat() would have returned an error
	// even though our big.Int will be fine. Thus, we let the error slide.

	num := big.Int{}
	if _, ok := num.SetString(s, 10); !ok {
		return nil, errors.New("invalid integer: " + s)
	}
	return &Number{true, f, *num}, nil
}

// SubtractNumbers subtracts two numbers and returns the difference.
func SubtractNumbers(n1, n2 *Number) *Number {
	i1, i2 := n1.Int(), n2.Int()
	if i1 != nil && i2 != nil {
		return NewNumberBig(big.NewInt(0).Sub(i1, i2))
	} else {
		return NewNumberFloat(n1.Float() - n2.Float())
	}
}
