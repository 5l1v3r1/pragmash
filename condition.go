package pragmash

type Condition []Runnable

func (c Condition) Run(r Runner) (Value, *Exception) {
	// TODO: this
	if len(c) == 0 {
	}
	return nil, nil
}

