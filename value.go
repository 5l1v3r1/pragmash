package pragmash

import (
	"bytes"
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var emptyValue = NewValueBool(false)

// A Value is a string which can also be treated as an array, a boolean and a number.
//
// Since all values in pragmash are strings, certain operations require converting a string to a
// different datatype. The Value type facilitates this conversion and caches the results of
// conversions.
//
// A Value is not thread-safe. If you intend to access it from multiple goroutines, you must
// synchronize such accesses.
type Value struct {
	arrayRep  []*Value
	boolRep   bool
	numRep    *Number
	numErr    error
	stringRep *string
}

// NewValueArray creates a new Value from an array.
func NewValueArray(arr []*Value) *Value {
	// If there is one empty element, the array must be empty in order to
	// maintain integrity.
	if len(arr) == 1 && len(arr[0].String()) == 0 {
		str := ""
		return &Value{[]*Value{}, false, nil, nil, &str}
	}
	return &Value{arr, len(arr) != 0, nil, nil, nil}
}

// NewValueBool creates a new Value from a boolean.
func NewValueBool(b bool) *Value {
	res := &Value{nil, b, nil, nil, nil}
	if b {
		str := "true"
		res.stringRep = &str
	} else {
		str := ""
		res.stringRep = &str
	}
	res.arrayRep = []*Value{res}
	return res
}

// NewValueString creates a new HybridValue from a string.
func NewValueString(str string) *Value {
	return &Value{nil, len(str) > 0, nil, nil, &str}
}

// NewValueNumber creates a new Value from a *Number.
func NewValueNumber(num *Number) *Value {
	res := &Value{nil, true, num, nil, nil}
	res.arrayRep = []*Value{res}
	return res
}

// Array returns an array which represents the value.
func (v *Value) Array() []*Value {
	if v.arrayRep != nil {
		return v.arrayRep
	}

	// Generate an array by splitting the string into parts.
	strVal := v.String()
	if len(strVal) == 0 {
		v.arrayRep = []*Value{}
		return v.arrayRep
	}
	comps := strings.Split(strVal, "\n")
	res := make([]*Value, len(comps))
	for i, x := range comps {
		res[i] = NewValueString(x)
	}
	v.arrayRep = res
	return res
}

// Bool returns the pre-cached boolean representation of the value.
func (v *Value) Bool() bool {
	return v.boolRep
}

// Number returns the numerical representation of the value, parsing it if necessary.
func (v *Value) Number() (*Number, error) {
	if v.numRep != nil || v.numErr != nil {
		return v.numRep, v.numErr
	}
	v.numRep, v.numErr = ParseNumber(v.String())
	return v.numRep, v.numErr
}

// String returns the string representation of the value.
func (v *Value) String() string {
	if v.stringRep != nil {
		return *v.stringRep
	} else if v.numRep != nil {
		str := v.numRep.String()
		v.stringRep = &str
		return str
	} else if v.arrayRep != nil {
		var buffer bytes.Buffer
		for i, elem := range v.arrayRep {
			if i != 0 {
				buffer.WriteRune('\n')
			}
			buffer.WriteString(elem.String())
		}
		str := buffer.String()
		v.stringRep = &str
		return str
	}
	panic("no way to generate a string representation")
	return ""
}

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
	res := Number{isInteger: true, floating: f}
	res.integer.Set(b)
	return &res
}

// NewNumberFloat generates a Number which probably will not have an integer
// representation.
func NewNumberFloat(f float64) *Number {
	rat := big.Rat{}
	rat.SetFloat64(f)
	if rat.IsInt() {
		return NewNumberBig(rat.Num())
	}
	return &Number{isInteger: false, floating: f}
}

// NewNumberInt generates a Number with an integer.
func NewNumberInt(i int64) *Number {
	return NewNumberBig(big.NewInt(i))
}

// Float returns the floating point representation of the number. All numbers have such a
// representation.
func (n *Number) Float() float64 {
	return n.floating
}

// Int returns the integer representation of the number if it exists.
// The Number's *big.Int is cached internally and is not copied before it is returned. As a result,
// you should not modify the *big.Int that you get from this method.
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

// String returns the string representation of the number.
func (n *Number) String() string {
	if n.isInteger {
		return n.integer.String()
	}
	return strconv.FormatFloat(n.floating, 'f', -1, 64)
}

// Zero returns true if the number is zero.
func (n *Number) Zero() bool {
	if n.isInteger {
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
		return NewNumberFloat(f), nil
	}

	// Parse it as a big int.

	// NOTE: if the number was HUGE, ParseFloat() would have returned an error
	// even though our big.Int will be fine. Thus, we let the error slide.

	num := big.Int{}
	if _, ok := num.SetString(s, 10); !ok {
		return nil, errors.New("invalid integer: " + s)
	}
	return &Number{true, f, num}, nil
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
