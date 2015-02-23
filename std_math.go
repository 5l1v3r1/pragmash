package pragmash

import (
	"math"
	"math/big"
	"math/rand"
)

// StdMath implements the standard arithmetic functions.
type StdMath struct{}

// Add adds a list of numbers.
func (_ StdMath) Add(nums []Number) Value {
	res := NewNumberInt(0)
	for _, num := range nums {
		res = AddNumbers(res, num)
	}
	return res
}

// Ceil returns the greatest integer which is less than or equal to a floating
// point.
func (_ StdMath) Ceil(num Number) Value {
	f := num.Float()
	rounded := math.Ceil(f)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewHybridValueNumber(NewNumberBig(rat.Num()))
}

// Div divides the first argument by the second.
func (_ StdMath) Div(n1, n2 Number) (Value, error) {
	return DivideNumbers(n1, n2)
}

// Floor returns the lowest integer which is greater than or equal to a floating
// point.
func (_ StdMath) Floor(num Number) Value {
	f := num.Float()
	rounded := math.Floor(f)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewHybridValueNumber(NewNumberBig(rat.Num()))
}

// Mul multiplies a list of numbers.
func (_ StdMath) Mul(nums []Number) Value {
	res := NewNumberInt(1)
	for _, num := range nums {
		res = MultiplyNumbers(res, num)
	}
	return res
}

// Pow raises the first argument to the power of the second.
func (_ StdMath) Pow(n1, n2 Number) Value {
	return ExponentiateNumber(n1, n2)
}

// Rand generates a random floating point number in [0.0, 1.0).
func (_ StdMath) Rand() Value {
	f := rand.Float64()
	return NewHybridValueNumber(NewNumberFloat(f))
}

// Round turns a floating point into a whole number by rounding it.
func (_ StdMath) Round(num Number) Value {
	f := num.Float()
	rounded := math.Floor(f + 0.5)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewHybridValueNumber(NewNumberBig(rat.Num()))
}

// Sub subtracts the second argument from the first.
func (_ StdMath) Sub(n1, n2 Number) Value {
	return SubtractNumbers(n1, n2)
}
