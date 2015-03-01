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
	StdTime
}

// NewStdRunner returns a Runner which implements the standard library.
func NewStdRunner(variables map[string]*Value) Runner {
	runner := NewReflectRunner(StdAll{}, OperatorRewrites)
	
	// Copy variables if necessary.
	if variables != nil {
		for name, value := range variables {
			runner.variables[name] = value
		}
	}
	
	return runner
}
