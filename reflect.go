package pragmash

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// A ReflectRunner implements a RunCommand() function that uses reflection.
type ReflectRunner struct {
	rewrite map[string]string
	value   reflect.Value
}

// NewReflectRunner creates a new ReflectRunner.
func NewReflectRunner(val interface{}, rw map[string]string) ReflectRunner {
	return ReflectRunner{rw, reflect.ValueOf(val)}
}

// RunCommand puts the name through the alias table if possible.
// It then capitalizes the first letter of the name and looks for a
// corresponding method.
func (r ReflectRunner) RunCommand(name string, v []Value) (Value, error) {
	n := r.RewriteName(name)
	n = strings.ToUpper(n[:1]) + n[1:]
	method := r.value.MethodByName(name)
	if !method.IsValid() {
		return nil, errors.New("unknown command: " + name)
	}
	t := method.Type()
	valType := reflect.TypeOf((*Value)(nil)).Elem()

	// Generate the arguments for the call.
	var args []reflect.Value
	if t.NumIn() == 0 {
		args = []reflect.Value{}
	} else if t.NumIn() == 1 && t.In(0) == reflect.TypeOf([]Value{}) {
		args = []reflect.Value{reflect.ValueOf(v)}
	} else if t.NumIn() != len(v) {
		return nil, errors.New("expected " + strconv.Itoa(t.NumIn()) +
			" arguments")
	} else {
		args = make([]reflect.Value, t.NumIn())
		for i, x := range v {
			if t.In(i) != valType {
				return nil, errors.New("invalid argument type")
			}
			args[i] = reflect.ValueOf(x)
		}
	}

	res := method.Call(args)
	if len(res) == 0 {
		return StringValue(""), nil
	} else if len(res) == 1 {
		// TODO: the result could be an error or a Value.
		return nil, nil
	} else if len(res) == 2 {
		// TODO: assert that the return types were (Value, error).
		return nil, nil
	}
	return nil, errors.New("invalid number of return values")
}

func (r ReflectRunner) RewriteName(name string) string {
	if r.rewrite != nil {
		if n, ok := r.rewrite[name]; ok {
			return n
		}
	}
	return name
}
