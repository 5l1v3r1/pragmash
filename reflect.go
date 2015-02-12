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
func (r ReflectRunner) RunCommand(name string, vals []Value) (Value, error) {
	n := r.RewriteName(name)
	n = strings.ToUpper(n[:1]) + n[1:]
	method := r.value.MethodByName(name)
	if !method.IsValid() {
		return nil, errors.New("unknown command: " + name)
	}
	t := method.Type()

	// Generate the arguments for the call.
	var args []reflect.Value
	if t.NumIn() == 0 {
		args = []reflect.Value{}
	} else if t.NumIn() == 1 && t.In(0) == reflect.TypeOf([]Value{}) {
		args = []reflect.Value{reflect.ValueOf(vals)}
	} else if t.NumIn() != len(vals) {
		return nil, errors.New("expected " + strconv.Itoa(t.NumIn()) +
			" arguments")
	} else {
		valType := reflect.TypeOf((*Value)(nil)).Elem()
		args = make([]reflect.Value, t.NumIn())
		for i, x := range vals {
			// TODO: here, support bool, string, and []string types.
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
		// The return type may be an error or a value.
		val := res[0].Interface()
		if err, ok := val.(error); ok {
			if err != nil {
				return nil, err
			} else {
				return StringValue(""), nil
			}
		} else if retVal, ok := val.(Value); ok {
			return retVal, nil
		}
		return nil, errors.New("invalid return type")
	} else if len(res) == 2 {
		// The return type must be (Value, error)
		i1 := res[0].Interface()
		i2 := res[0].Interface()
		if val, ok := i1.(Value); !ok {
			return nil, errors.New("invalid first return type")
		} else if err, ok := i2.(error); !ok {
			return nil, errors.New("invalid second return type")
		} else {
			return val, err
		}
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
