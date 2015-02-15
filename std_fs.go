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
func (_ StdFs) Exists(path string) (Value, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return BoolValue(false), nil
	} else if err != nil {
		return nil, err
	}
	return BoolValue(true), nil
}

// Glob matches filenames with wildcards.
func (_ StdFs) Glob(args []Value) (Value, error) {
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

// Mkdir creates a directory or fails with an error.
func (_ StdFs) Mkdir(name string) error {
	return os.Mkdir(name, os.FileMode(0755))
}

// Path joins path components.
func (_ StdFs) Path(args []Value) Value {
	comps := make([]string, len(args))
	for i, x := range args {
		comps[i] = x.String()
	}
	return StringValue(filepath.Join(comps...))
}

// Rm removes a file or directory but does not do so recursively.
func (_ StdFs) Rm(path string) error {
	return os.Remove(path)
}

// Rmall removes a file or directory and all its children.
func (_ StdFs) Rmall(path string) error {
	return os.RemoveAll(path)
}
