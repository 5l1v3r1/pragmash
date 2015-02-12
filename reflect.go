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
	// Lookup the method.
	n := r.RewriteName(name)
	n = strings.ToUpper(n[:1]) + n[1:]
	method := r.value.MethodByName(name)
	if !method.IsValid() {
		return nil, errors.New("unknown command: " + name)
	}
	t := method.Type()

	// Generate the arguments.
	args, err := reflectArguments(t, vals)
	if err != nil {
		return nil, err
	}

	// Run the call and process the return value.
	res := method.Call(args)
	return reflectReturnValue(res)
}

func (r ReflectRunner) RewriteName(name string) string {
	if r.rewrite != nil {
		if n, ok := r.rewrite[name]; ok {
			return n
		}
	}
	return name
}

func reflectArguments(t reflect.Type, vals []Value) ([]reflect.Value, error) {
	// Special cases.
	if t.NumIn() == 0 {
		return []reflect.Value{}, nil
	} else if t.NumIn() == 1 && t.In(0) == reflect.TypeOf([]Number{}) {
		// Generate a list of numbers.
		nums := make([]Number, len(vals))
		for i, x := range vals {
			num, err := x.Number()
			if err != nil {
				return nil, err
			}
			nums[i] = num
		}
		return []reflect.Value{reflect.ValueOf(nums)}, nil
	} else if t.NumIn() == 1 && t.In(0) == reflect.TypeOf([]Value{}) {
		return []reflect.Value{reflect.ValueOf(vals)}, nil
	} else if t.NumIn() != len(vals) {
		return nil, errors.New("expected " + strconv.Itoa(t.NumIn()) +
			" arguments")
	}

	// These are the allowed argument types.
	arrType := reflect.TypeOf([]string{})
	boolType := reflect.TypeOf(true)
	numType := reflect.TypeOf((*Number)(nil)).Elem()
	strType := reflect.TypeOf("")
	valType := reflect.TypeOf((*Value)(nil)).Elem()

	// Process each argument individually.
	args := make([]reflect.Value, t.NumIn())
	for i, x := range vals {
		inputType := t.In(i)
		if inputType == valType {
			args[i] = reflect.ValueOf(x)
		} else if inputType == numType {
			num, err := x.Number()
			if err != nil {
				return nil, err
			}
			args[i] = reflect.ValueOf(num)
		} else if inputType == boolType {
			args[i] = reflect.ValueOf(x.Bool())
		} else if inputType == arrType {
			args[i] = reflect.ValueOf(x.Array())
		} else if inputType == strType {
			args[i] = reflect.ValueOf(x.String())
		} else {
			return nil, errors.New("invalid argument type: " +
				inputType.String())
		}
	}

	return args, nil
}

func reflectReturnValue(res []reflect.Value) (Value, error) {
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
