package pragmash

// CommandRunnable is a runnable which executes a command.
type CommandRunnable struct {
	Context   string
	Name      Runnable
	Arguments []Runnable
}

// Run evaluates every argument, then executes the named command.
func (c CommandRunnable) Run(r Runner) (Value, *Exception) {
	name, exc := c.Name.Run(r)
	if exc != nil {
		return nil, exc
	}
	args := make([]Value, len(c.Arguments))
	for i, x := range c.Arguments {
		val, exc := x.Run(r)
		if exc != nil {
			return nil, exc
		}
		args[i] = val
	}
	val, err := r.RunCommand(name.String(), args)
	if err != nil {
		return nil, NewException(c.Context, err)
	}
	return val, nil
}

// A Runnable is a generic interface which can execute on a given Runner and
// return a value or runtime exception.
type Runnable interface {
	Run(r Runner) (Value, *Exception)
}

// A RunnableList is a Runnable which runs a list of Runnables.
type RunnableList []Runnable

// Run runs each Runnable and fails on the first exception it encounters.
// If no exception is encountered, this returns an empty string.
func (r RunnableList) Run(runner Runner) (Value, *Exception) {
	for _, x := range r {
		if _, exc := x.Run(runner); exc != nil {
			return nil, exc
		}
	}
	return StringValue(""), nil
}

// A Runner is a generic interface which can run a commands.
type Runner interface {
	RunCommand(name string, args []Value) (Value, error)
}

// A ValueRunnable always returns the same value.
type ValueRunnable struct {
	Value
}

// Run returns Value(v), nil.
func (v ValueRunnable) Run(r Runner) (Value, *Exception) {
	return v.Value, nil
}
