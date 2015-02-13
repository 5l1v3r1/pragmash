package pragmash

// StdMath implements the standard arithmetic functions.
type StdMath struct{}

// Add adds a list of numbers.
func (s StdMath) Add(nums []Number) Value {
	res := NewNumberInt(0)
	for _, num := range nums {
		res = AddNumbers(res, num)
	}
	return res
}

// Div divides the first argument by the second.
func (s StdMath) Div(n1, n2 Number) (Value, error) {
	return DivideNumbers(n1, n2)
}

// Mul multiplies a list of numbers.
func (s StdMath) Mul(nums []Number) Value {
	res := NewNumberInt(1)
	for _, num := range nums {
		res = MultiplyNumbers(res, num)
	}
	return res
}

// Pow raises the first argument to the power of the second.
func (s StdMath) Pow(n1, n2 Number) Value {
	return ExponentiateNumber(n1, n2)
}

// Sub subtracts the second argument from the first.
func (s StdMath) Sub(n1, n2 Number) Value {
	return SubtractNumbers(n1, n2)
}
