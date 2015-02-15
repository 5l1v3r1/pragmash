package pragmash

// StdAll implements the methods corresponding to the standard library.
type StdAll struct {
	StdArray
	StdFs
	StdInternal
	StdIo
	StdMath
	StdOps
	StdString
}

// NewStdAll creates a new StdAll instance with a new StdInternal inside.
func NewStdAll() StdAll {
	return StdAll{StdInternal: NewStdInternal()}
}

// NewStdRunner returns a Runner which implements the standard library.
func NewStdRunner(variables map[string]Value) Runner {
	var runner Runner
	all := NewStdAll()
	all.Runner = &runner
	if variables != nil {
		for variable, value := range variables {
			all.Variables[variable] = value
		}
	}
	res := NewReflectRunner(all, OperatorRewrites)
	runner = res
	return res
}
