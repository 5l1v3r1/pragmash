package pragmash

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// StdFs provides commands for file system manipulation.
type StdFs struct{}

// Exists returns whether or not a file exists.
func (s StdFs) Exists(path string) (Value, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return BoolValue(false), nil
	} else if err != nil {
		return nil, err
	}
	return BoolValue(true), nil
}

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

// Rm removes a file or directory but does not do so recursively.
func (s StdFs) Rm(path string) error {
	return os.Remove(path)
}

// Rmall removes a file or directory and all its children.
func (s StdFs) Rmall(path string) error {
	return os.RemoveAll(path)
}
