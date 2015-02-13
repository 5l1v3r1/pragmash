package pragmash

// A Condition is a boolean expression used in "if" statements and "while"
// loops.
type Condition []Runnable

// Run evaluates the condition and returns a BoolValue on success.
// Empty conditions are automatically true. Conditions with one argument are
// true if the argument's length is non-zero. Conditions with more than one
// argument are true if all the arguments equal.
func (c Condition) Run(r Runner) (Value, *Exception) {
	if len(c) == 0 {
		// Empty conditions are automatically true.
		return BoolValue(true), nil
	} else if len(c) == 1 {
		// Every non-empty string is true.
		val, exc := c[0].Run(r)
		if exc != nil {
			return nil, exc
		}
		return BoolValue(val.Bool()), nil
	}

	// Make sure every value equals the first.
	first, exc := c[0].Run(r)
	if exc != nil {
		return nil, exc
	}
	str := first.String()
	for i := 1; i < len(c); i++ {
		val, exc := c[i].Run(r)
		if exc != nil {
			return nil, exc
		}
		if val.String() != str {
			return BoolValue(false), nil
		}
	}

	return BoolValue(true), nil
}

// A NotCondition is essentially the inverse of a Condition.
type NotCondition []Runnable

// Run evaluates the NotCondition and returns a BoolValue on success.
// Empty conditions are automatically false. Conditions with one argument are
// true if the argument's length is zero. Conditions with more than one
// argument are true if at least one of the arguments differs.
func (n NotCondition) Run(r Runner) (Value, *Exception) {
	val, exc := Condition(n).Run(r)
	if exc != nil {
		return nil, exc
	}
	return BoolValue(!val.Bool()), nil
}

// ConditionFromTokens reads a series of tokens and converts them into a
// Runnable which is either a Condition or a NotCondition.
func ConditionFromTokens(t []Token, context string) Runnable {
	if len(t) != 0 && t[0].String == "not" {
		// Negative condition
		c := make(NotCondition, len(t)-1)
		for i := 1; i < len(t); i++ {
			c[i-1] = t[i].Runnable(context)
		}
		return c
	} else {
		// Positive condition
		c := make(Condition, len(t))
		for i, token := range t {
			c[i] = token.Runnable(context)
		}
		return c
	}
}
