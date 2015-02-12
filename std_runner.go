package pragmash

// StdAll implements the methods corresponding to the standard library.
type StdAll struct {
	StdGenerators
	StdInternal
	StdIo
	StdMath
	StdOps
}

// NewStdRunner returns a Runner which implements the standard library.
func NewStdRunner() Runner {
	return NewReflectRunner(StdAll{}, OperatorRewrites)
}
