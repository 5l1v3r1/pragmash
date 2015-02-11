package pragmash

// An If object represents an "if" statement.
type If struct {
	Branches   []Runnable
	Conditions []Runnable
}

// Run executes the if statement.
func (block *If) Run(r Runner) (Value, *Exception) {
	for i, cond := range block.Conditions {
		val, exc := cond.Run(r)
		if exc != nil {
			return nil, exc
		}
		if val.Bool() {
			// Run the branch.
			return block.Branches[i].Run(r)
		}
	}
	// No branches ran.
	return StringValue(""), nil
}
