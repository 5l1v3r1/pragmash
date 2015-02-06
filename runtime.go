package pragmash

// A Runnable is a generic interface which can execute on a given Runner and
// return a value or runtime exception.
type Runnable interface {
	Run(r Runner) (string, *Exception)
}

// A Runner is a generic interface which can run a commands.
type Runner interface {
	RunCommand(name string, args []Value) (Value, error)
}
