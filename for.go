package pragmash

// A For block represents a for-loop.
type For struct {
	Body       Runnable
	Context    string
	Expression Runnable
	Variable   *Runnable
}

// Run executes the for loop.
// This fails if the variable name, exression, or body triggers an exception.
func (f For) Run(r Runner) (Value, *Exception) {
	expr, exc := f.Expression.Run(r)
	if exc != nil {
		return nil, exc
	}
	var variable Value
	if f.Variable != nil {
		variable, exc = (*f.Variable).Run(r)
		if exc != nil {
			return nil, exc
		}
	}
	for _, val := range expr.Array() {
		if variable != nil {
			_, err := r.RunCommand("set", []Value{variable, val})
			if err != nil {
				return nil, NewException(f.Context, err)
			}
		}
		_, exc = f.Body.Run(r)
		if exc != nil {
			return nil, exc
		}
	}
	return StringValue(""), nil
}
