package pragmash

import (
	"errors"
	"strconv"
)

// OperatorRewrites can be used for a ReflectRunner's command rewrite table to
// replace the named operators with symbolic ones like *, +, [], <, etc.
var OperatorRewrites = map[string]string{
	"+": "add", "/": "div", "*": "mul", "-": "sub", "**": "pow", "%": "mod",
	"[]": "subscript", "<=": "le", ">=": "ge", "<": "lt", ">": "gt", "=": "eq",
	"||": "or", "&&": "and",
}

// StdOps implements the standard operators that are not implemented in StdMath.
type StdOps struct{}

// And implements the && operator.
func (_ StdOps) And(vals []*Value) *Value {
	for _, v := range vals {
		if !v.Bool() {
			return emptyValue
		}
	}
	return NewValueBool(true)
}

// Eq implements the equality operator.
func (_ StdOps) Eq(s1, s2 string) *Value {
	return NewValueBool(s1 == s2)
}

// Ge implements the >= operator.
func (_ StdOps) Ge(n1, n2 *Number) *Value {
	return NewValueBool(CompareNumbers(n1, n2) >= 0)
}

// Gt implements the > operator.
func (_ StdOps) Gt(n1, n2 *Number) *Value {
	return NewValueBool(CompareNumbers(n1, n2) > 0)
}

// Le implements the <= operator.
func (_ StdOps) Le(n1, n2 *Number) *Value {
	return NewValueBool(CompareNumbers(n1, n2) <= 0)
}

// Lt implements the < operator.
func (_ StdOps) Lt(n1, n2 *Number) *Value {
	return NewValueBool(CompareNumbers(n1, n2) < 0)
}

// Or implements the || operator.
func (_ StdOps) Or(vals []*Value) *Value {
	for _, v := range vals {
		if v.Bool() {
			return v
		}
	}
	return emptyValue
}

// Subscript gets a term from a list.
func (_ StdOps) Subscript(vals []*Value, index int) (*Value, error) {
	if index < 0 || index >= len(vals) {
		return nil, errors.New("subscript out of bounds: " +
			strconv.Itoa(index))
	}
	return vals[index], nil
}
