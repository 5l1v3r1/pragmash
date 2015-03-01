package pragmash

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestAllScripts(t *testing.T) {
	listing, err := listTestDirectory()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range listing {
		testName := filepath.Base(path)
		s, err := readTestScript(path)
		if err != nil {
			t.Error(err)
			continue
		}
		if err := s.run(); err != nil {
			t.Error("error in " + testName + ": " + err.Error())
		}
	}
}

type testScript struct {
	expect string
	path   string
}

func (t *testScript) run() error {
	variables := map[string]*Value{
		"ARGV": NewValueArray([]*Value{}),
		"DIR": NewValueString(filepath.Dir(t.path)),
	}
	runner := NewStdRunner(variables)
	
	// Create the script
	contents, err := ioutil.ReadFile(t.path)
	if err != nil {
		return err
	}
	
	// Parse the script
	lines, contexts, err := TokenizeString(string(contents))
	if err != nil {
		return err
	}
	runnable, err := ScanAll(lines, contexts)
	if err != nil {
		return err
	}
	
	if _, bo := runnable.Run(runner); bo == nil {
		return errors.New("no breakout")
	} else if bo.Type() != BreakoutTypeReturn {
		return errors.New("unexpected breakout")
	} else if bo.Value().String() != t.expect {
		return errors.New("unexpected output: " + bo.Value().String())
	}
	return nil
}

func listTestDirectory() ([]string, error) {
	_, filename, _, _ := runtime.Caller(0)
	testsPath := filepath.Join(filepath.Dir(filename), "tests")
	f, err := os.Open(testsPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(names))
	for _, x := range names {
		if !strings.HasSuffix(x, ".pragmash") {
			continue
		}
		res = append(res, filepath.Join(testsPath, x))
	}
	return res, nil
}

func readTestScript(path string) (*testScript, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contents := string(data)

	// Read the first line which contains the quoted string for the expected
	// output.
	lines := strings.Split(contents, "\n")
	first := lines[0]
	if len(first) < 4 {
		return nil, errors.New("invalid first line")
	}
	if first[:3] != "# \"" {
		return nil, errors.New("invalid first line")
	}
	scanner := NewScannerString(first[3:])
	expect, err := scanner.ReadQuoted()
	if err != nil {
		return nil, err
	}

	return &testScript{expect, path}, nil
}
