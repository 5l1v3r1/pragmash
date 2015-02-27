package pragmash

// CommandRunnable is a runnable which executes a command.
type CommandRunnable struct {
	Context   string
	Name      Runnable
	Arguments []Runnable
}

// Run evaluates every argument, then executes the named command.
func (c CommandRunnable) Run(r Runner) (*Value, *Breakout) {
	name, exc := c.Name.Run(r)
	if exc != nil {
		return nil, exc
	}
	args := make([]*Value, len(c.Arguments))
	for i, x := range c.Arguments {
		val, bo := x.Run(r)
		if bo != nil {
			return nil, bo
		}
		args[i] = val
	}
	val, err := r.RunCommand(name.String(), args)
	if err != nil {
		return nil, NewBreakoutException(c.Context, err)
	}
	return val, nil
}

// A Runnable is a generic interface which can execute on a given Runner and
// return a value or breakout.
type Runnable interface {
	Run(r Runner) (*Value, *Breakout)
}

// A RunnableList is a Runnable which runs a list of Runnables.
type RunnableList []Runnable

// Run runs each Runnable and fails on the first breakout it encounters.
// If no breakout is encountered, this returns the value of the last runnable.
func (r RunnableList) Run(runner Runner) (*Value, *Breakout) {
	lastValue := emptyValue
	for _, x := range r {
		if val, bo := x.Run(runner); bo != nil {
			return nil, bo
		} else {
			lastValue = val
		}
	}
	return lastValue, nil
}

// A Runner is a generic interface which can run a commands.
type Runner interface {
	RunCommand(name string, args []*Value) (*Value, error)
}
