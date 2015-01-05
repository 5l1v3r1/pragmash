package pragmash

// Context is a context which can evaluate commands.
type Context interface {
	// Evaluate evaluates a flattened command.
	Evaluate(c *Command) (string, error)
}
