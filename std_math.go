package pragmash

import (
	"math"
	"math/big"
	"math/rand"
)

// StdMath implements the standard arithmetic functions.
type StdMath struct{}

// Abs returns the absolute value of a number.
func (_ StdMath) Abs(num *Number) *Value {
	// TODO: take a Value argument so we can possibly return the same *Value...
	if CompareNumbers(num, NewNumberInt(0)) == -1 {
		res := MultiplyNumbers(NewNumberInt(-1), num)
		return NewValueNumber(res)
	}
	return NewValueNumber(num)
}

// Add adds a list of numbers.
func (_ StdMath) Add(nums []*Number) *Value {
	res := NewNumberInt(0)
	for _, num := range nums {
		res = AddNumbers(res, num)
	}
	return NewValueNumber(res)
}

// Ceil returns the greatest integer which is less than or equal to a floating
// point.
func (_ StdMath) Ceil(num *Number) *Value {
	f := num.Float()
	rounded := math.Ceil(f)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewValueNumber(NewNumberBig(rat.Num()))
}

// Cos returns the cosine of its argument (which is in radians).
func (_ StdMath) Cos(num *Number) *Value {
	res := math.Cos(num.Float())
	return NewValueNumber(NewNumberFloat(res))
}

// Div divides the first argument by the second.
func (_ StdMath) Div(n1, n2 *Number) (*Value, error) {
	num, err := DivideNumbers(n1, n2)
	if err != nil {
		return nil, err
	}
	return NewValueNumber(num), nil
}

// Floor returns the lowest integer which is greater than or equal to a floating
// point.
func (_ StdMath) Floor(num *Number) *Value {
	f := num.Float()
	rounded := math.Floor(f)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewValueNumber(NewNumberBig(rat.Num()))
}

// Mod computes the remainder of a division operation.
func (_ StdMath) Mod(num, modulus *Number) *Value {
	i1, i2 := num.Int(), modulus.Int()
	if i1 == nil || i2 == nil {
		// Do a funky floating-point modulus.
		// Some languages like Processing do this; I figure I might as well.
		f1, f2 := num.Float(), modulus.Float()
		quot := math.Floor(f1 / f2)
		res := f1 - quot*f2
		return NewValueNumber(NewNumberFloat(res))
	}
	var resNum big.Int
	resNum.Mod(i1, i2)
	return NewValueNumber(NewNumberBig(&resNum))
}

// Mul multiplies a list of numbers.
func (_ StdMath) Mul(nums []*Number) *Value {
	res := NewNumberInt(1)
	for _, num := range nums {
		res = MultiplyNumbers(res, num)
	}
	return NewValueNumber(res)
}

// Pi returns the value of pi.
func (_ StdMath) Pi() *Value {
	return NewValueNumber(NewNumberFloat(math.Pi))
}

// Pow raises the first argument to the power of the second.
func (_ StdMath) Pow(n1, n2 *Number) *Value {
	return NewValueNumber(ExponentiateNumber(n1, n2))
}

// Rand generates a random floating point number in [0.0, 1.0).
func (_ StdMath) Rand() *Value {
	f := rand.Float64()
	return NewValueNumber(NewNumberFloat(f))
}

// Round turns a floating point into a whole number by rounding it.
func (_ StdMath) Round(num *Number) *Value {
	f := num.Float()
	rounded := math.Floor(f + 0.5)
	rat := big.NewRat(0, 1)
	rat.SetFloat64(rounded)
	return NewValueNumber(NewNumberBig(rat.Num()))
}

// Sin returns the sine of its argument (which is in radians).
func (_ StdMath) Sin(num *Number) *Value {
	res := math.Sin(num.Float())
	return NewValueNumber(NewNumberFloat(res))
}

// Sub subtracts the second argument from the first.
func (_ StdMath) Sub(n1, n2 *Number) *Value {
	return NewValueNumber(SubtractNumbers(n1, n2))
}
