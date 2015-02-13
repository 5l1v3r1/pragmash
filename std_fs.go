package pragmash

import (
	"path/filepath"
	"sort"
	"strings"
)

// StdFs provides commands for file system manipulation.
type StdFs struct{}

// Glob matches filenames with wildcards.
func (s StdFs) Glob(args []Value) (Value, error) {
	res := make([]string, 0)
	for _, v := range args {
		paths, err := filepath.Glob(v.String())
		if err != nil {
			return nil, err
		}
		res = append(res, paths...)
	}
	sort.Strings(res)
	return StringValue(strings.Join(res, "\n")), nil
}
