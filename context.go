package pragmash

type CommandFunc func([]string) (string, error)

type Context struct {
	Variables map[string]string
	Commands  map[string]CommandFunc
}

type Expression interface {
	Run(c *Context) (string, error)
}

