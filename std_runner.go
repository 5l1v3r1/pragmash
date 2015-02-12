package pragmash

// StdAll implements the methods corresponding to the standard library.
type StdAll struct {
	StdGenerators
	StdInternal
	StdIo
	StdMath
	StdOps
}

// NewStdAll creates a new StdAll instance with a new StdInternal inside.
func NewStdAll() StdAll {
	return StdAll{StdInternal: NewStdInternal()}
}

// NewStdRunner returns a Runner which implements the standard library.
func NewStdRunner() Runner {
	return NewReflectRunner(NewStdAll(), OperatorRewrites)
}
