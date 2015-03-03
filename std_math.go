package pragmash

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
)

// StdMath implements the standard arithmetic functions.
type StdMath struct{}

// Abs returns the absolute value of a number.
func (_ StdMath) Abs(num *Number) *Number {
	// TODO: take a Value argument so we can possibly return the same *Value...
	if CompareNumbers(num, NewNumberInt(0)) == -1 {
		return MultiplyNumbers(NewNumberInt(-1), num)
	}
	return num
}

// Add adds a list of numbers.
func (_ StdMath) Add(nums ...*Number) *Number {
	res := NewNumberInt(0)
	for _, num := range nums {
		res = AddNumbers(res, num)
	}
	return res
}

// Ceil returns the greatest integer which is less than or equal to a floating
// point.
func (_ StdMath) Ceil(f float64) *Number {
	rounded := math.Ceil(f)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewNumberBig(rat.Num())
}

// Cos returns the cosine of its argument (which is in radians).
func (_ StdMath) Cos(f float64) float64 {
	return math.Cos(f)
}

// Div divides the first argument by the second.
func (_ StdMath) Div(n1, n2 *Number) (*Number, error) {
	return DivideNumbers(n1, n2)
}

// Floor returns the lowest integer which is greater than or equal to a floating
// point.
func (_ StdMath) Floor(f float64) *Number {
	rounded := math.Floor(f)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewNumberBig(rat.Num())
}

// Log computes a logarithm.
// If Log is given one argument, this computes log base 10.
// If two arguments are given, the first is used as the base and the second is
// used as the argument.
func (s StdMath) Log(arg1 float64, args ...float64) (float64, error) {
	if len(args) > 1 {
		return 0, errors.New("expected 1 or 2 arguments")
	}
	if len(args) == 0 {
		return s.Log(10, arg1)
	}
	conversion := math.Log(arg1)
	if math.IsNaN(conversion) || math.IsInf(conversion, 0) || conversion == 0 {
		return 0, errors.New("invalid base")
	}
	res := math.Log(args[0])
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, errors.New("invalid argument")
	}
	return res / conversion, nil
}

// Mod computes the remainder of a division operation.
func (_ StdMath) Mod(num, modulus *Number) *Number {
	i1, i2 := num.Int(), modulus.Int()
	if i1 == nil || i2 == nil {
		// Do a funky floating-point modulus.
		// Some languages like Processing do this; I figure I might as well.
		f1, f2 := num.Float(), modulus.Float()
		quot := math.Floor(f1 / f2)
		res := f1 - quot*f2
		return NewNumberFloat(res)
	}
	var resNum big.Int
	resNum.Mod(i1, i2)
	return NewNumberBig(&resNum)
}

// Mul multiplies a list of numbers.
func (_ StdMath) Mul(nums ...*Number) *Number {
	res := NewNumberInt(1)
	for _, num := range nums {
		res = MultiplyNumbers(res, num)
	}
	return res
}

// Pi returns the value of pi.
func (_ StdMath) Pi() float64 {
	return math.Pi
}

// Pow raises the first argument to the power of the second.
func (_ StdMath) Pow(n1, n2 *Number) *Number {
	return ExponentiateNumber(n1, n2)
}

// Rand generates a random floating point number in [0.0, 1.0).
func (_ StdMath) Rand() float64 {
	return rand.Float64()
}

// Round turns a floating point into a whole number by rounding it.
func (_ StdMath) Round(f float64) *Number {
	rounded := math.Floor(f + 0.5)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewNumberBig(rat.Num())
}

// Sin returns the sine of its argument (which is in radians).
func (_ StdMath) Sin(f float64) float64 {
	return math.Sin(f)
}

// Sub subtracts the second argument from the first.
func (_ StdMath) Sub(n1, n2 *Number) *Number {
	return SubtractNumbers(n1, n2)
}
