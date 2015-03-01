package pragmash

import (
	"os"
	"path/filepath"
	"sort"
)

// StdFs provides commands for file system manipulation.
type StdFs struct{}

// Exists returns whether or not a file exists.
func (_ StdFs) Exists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Filetype returns "file", "dir", "link", or "other" for a given file.
func (_ StdFs) Filetype(path string) (string, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	t := stat.Mode() & os.ModeType
	if t == 0 {
		return "file", nil
	} else if (t & os.ModeSymlink) != 0 {
		return "link", nil
	} else if (t & os.ModeDir) != 0 {
		return "dir", nil
	}
	return "other", nil
}

// Glob matches filenames with wildcards.
func (_ StdFs) Glob(args ...string) ([]string, error) {
	res := make([]string, 0)
	for _, name := range args {
		paths, err := filepath.Glob(name)
		if err != nil {
			return nil, err
		}
		res = append(res, paths...)
	}
	sort.Strings(res)
	return res, nil
}

// Mkdir creates a directory or fails with an error.
func (_ StdFs) Mkdir(name string) error {
	return os.Mkdir(name, os.FileMode(0755))
}

// Path joins path components.
func (_ StdFs) Path(comps ...string) string {
	return filepath.Join(comps...)
}

// Rm removes a file or directory but does not do so recursively.
func (_ StdFs) Rm(path string) error {
	return os.Remove(path)
}

// Rmall removes a file or directory and all its children.
func (_ StdFs) Rmall(path string) error {
	return os.RemoveAll(path)
}
