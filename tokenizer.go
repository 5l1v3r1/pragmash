package pragmash

type Line struct {
	CloseBlock bool
	OpenBlock  bool
	Tokens     []string
}

func TokenizeLine(line string)
