package pragmash

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// These are the allowed argument types.
var (
	boolType     = reflect.TypeOf(true)
	errType      = reflect.TypeOf((*error)(nil)).Elem()
	floatArrType = reflect.TypeOf([]float64{})
	floatType    = floatArrType.Elem()
	intArrType   = reflect.TypeOf([]int{})
	intType      = intArrType.Elem()
	numArrType   = reflect.TypeOf([]*Number{})
	numType      = numArrType.Elem()
	runnerType   = reflect.TypeOf((*Runner)(nil)).Elem()
	strArrType   = reflect.TypeOf([]string{})
	strType      = strArrType.Elem()
	valArrType   = reflect.TypeOf([]*Value{})
	valType      = valArrType.Elem()
)

// A ReflectRunner implements a RunCommand() function that uses reflection.
type ReflectRunner struct {
	rewrite   map[string]string
	value     reflect.Value
	variables map[string]*Value
}

// NewReflectRunner creates a new ReflectRunner.
func NewReflectRunner(val interface{}, rw map[string]string) *ReflectRunner {
	return &ReflectRunner{rw, reflect.ValueOf(val), map[string]*Value{}}
}

// RunCommand puts the name through the alias table if possible.
// It then capitalizes the first letter of the name and looks for a
// corresponding method.
// This will execute a special subroutine for the set and get commands.
func (r *ReflectRunner) RunCommand(name string, vals []*Value) (*Value, error) {
	if name == "get" {
		return r.getCommand(vals)
	} else if name == "set" {
		return r.setCommand(vals)
	}

	// Lookup the method.
	n := r.RewriteName(name)
	method := r.value.MethodByName(n)
	if !method.IsValid() {
		return nil, errors.New("unknown command: " + name)
	}
	t := method.Type()

	// Generate the arguments.
	args, err := r.arguments(t, vals)
	if err != nil {
		return nil, err
	}

	// Run the call and process the return value.
	res := method.Call(args)
	return reflectReturnValue(res)
}

// RewriteName uses the ReflectRunner's rewrite table to rewrite a given command
// name. If no rewrite rule is found, underscores are replaced with camel case.
func (r *ReflectRunner) RewriteName(name string) string {
	if r.rewrite != nil {
		if n, ok := r.rewrite[name]; ok {
			name = n
		}
	}
	// Capitalize the first letter.
	name = strings.ToUpper(name[:1]) + name[1:]
	// Replace "a_b" with "aB"
	for i := 1; i < len(name)-1; i++ {
		if name[i] == '_' {
			name = name[:i] + strings.ToUpper(name[i+1 : i+2]) + name[i+2:]
		}
	}
	return name
}

func (r *ReflectRunner) arguments(t reflect.Type,
	vals []*Value) ([]reflect.Value, error) {
	// The resulting arguments will be appended to this slice.
	res := make([]reflect.Value, 0, t.NumIn())

	// This will be incremented whenever a value from vals is used.
	valIdx := 0

	// If there are variadic arguments, the last argument is actually a slice
	// and doesn't count towards the expected argument count.
	expectedArgs := t.NumIn()
	printArgs := expectedArgs
	if t.IsVariadic() {
		expectedArgs--
	}

	// Process the normal (non-variadic) arguments.
	for i := 0; i < expectedArgs; i++ {
		argType := t.In(i)

		// If the argument is a Runner, no value is associated with it.
		if argType == runnerType {
			res = append(res, reflect.ValueOf(r))
			printArgs--
			continue
		} else if valIdx == len(vals) {
			// They are missing at least one argument.
			return nil, argumentsError(t.IsVariadic(), printArgs)
		}

		// Process a regular argument.
		val, err := pragmashValueToGo(argType, vals[valIdx])
		if err != nil {
			return nil, err
		}
		res = append(res, val)
		valIdx++
	}

	// Process the variadic arguments if there are any.
	if t.IsVariadic() {
		argType := t.In(t.NumIn() - 1).Elem()
		for valIdx < len(vals) {
			val, err := pragmashValueToGo(argType, vals[valIdx])
			if err != nil {
				return nil, err
			}
			res = append(res, val)
			valIdx++
		}
	} else if valIdx < len(vals) {
		return nil, argumentsError(false, printArgs)
	}

	return res, nil
}

func (r *ReflectRunner) getCommand(vals []*Value) (*Value, error) {
	if len(vals) != 1 {
		return nil, errors.New("expected 1 argument")
	}
	name := vals[0].String()
	if v, ok := r.variables[name]; ok {
		return v, nil
	} else {
		return nil, errors.New("variable undefined: " + name)
	}
}

func (r *ReflectRunner) setCommand(vals []*Value) (*Value, error) {
	if len(vals) != 2 {
		return nil, errors.New("expected 2 arguments")
	}
	r.variables[vals[0].String()] = vals[1]
	return emptyValue, nil
}

func argumentsError(variadic bool, count int) error {
	// If it's variadic, we add "at least" to the error message.
	if variadic {
		if count == 1 {
			return errors.New("expected at least 1 argument")
		}
		return errors.New("expected at least " + strconv.Itoa(count) +
			" arguments")
	}

	if count == 1 {
		return errors.New("expected 1 argument")
	}
	return errors.New("expected " + strconv.Itoa(count) + " arguments")
}

func goValueToPragmash(v interface{}) (*Value, error) {
	switch v := v.(type) {
	case bool:
		return NewValueBool(v), nil
	case []float64:
		numbers := make([]*Value, len(v))
		for i, f := range v {
			numbers[i] = NewValueNumber(NewNumberFloat(f))
		}
		return NewValueArray(numbers), nil
	case float64:
		return NewValueNumber(NewNumberFloat(v)), nil
	case []int:
		numbers := make([]*Value, len(v))
		for i, x := range v {
			numbers[i] = NewValueNumber(NewNumberInt(int64(x)))
		}
		return NewValueArray(numbers), nil
	case int:
		return NewValueNumber(NewNumberInt(int64(v))), nil
	case []*Number:
		numbers := make([]*Value, len(v))
		for i, x := range v {
			numbers[i] = NewValueNumber(x)
		}
		return NewValueArray(numbers), nil
	case *Number:
		return NewValueNumber(v), nil
	case []string:
		values := make([]*Value, len(v))
		for i, s := range v {
			values[i] = NewValueString(s)
		}
		return NewValueArray(values), nil
	case string:
		return NewValueString(v), nil
	case []*Value:
		return NewValueArray(v), nil
	case *Value:
		return v, nil
	default:
		return nil, errors.New("unexpected return type")
	}
}

func pragmashValueToGo(t reflect.Type, v *Value) (reflect.Value, error) {
	switch t {
	case boolType:
		return reflect.ValueOf(v.Bool()), nil
	case floatArrType:
		return valueToFloatArray(v)
	case floatType:
		return valueToFloat(v)
	case intArrType:
		return valueToIntArray(v)
	case intType:
		return valueToInt(v)
	case numArrType:
		return valueToNumArray(v)
	case numType:
		return valueToNum(v)
	case strArrType:
		return valueToStrArray(v)
	case strType:
		return reflect.ValueOf(v.String()), nil
	case valArrType:
		return reflect.ValueOf(v.Array()), nil
	case valType:
		return reflect.ValueOf(v), nil
	default:
		return reflect.ValueOf(nil), errors.New("unknown argument type")
	}
}

func reflectReturnValue(res []reflect.Value) (*Value, error) {
	// If there was no return value, this is easy.
	if len(res) == 0 {
		return emptyValue, nil
	}

	// If there was a single return value, it might be an error.
	if len(res) == 1 {
		if res[0].Type() == errType {
			val := res[0].Interface()
			if val != nil {
				return nil, val.(error)
			} else {
				return emptyValue, nil
			}
		}
		return goValueToPragmash(res[0].Interface())
	}

	// The return type must be (SOMEVALUE, error)
	if len(res) != 2 {
		return nil, errors.New("invalid number of return values")
	} else if res[1].Type() != errType {
		return nil, errors.New("invalid second return type")
	}
	if errVal := res[1].Interface(); errVal != nil {
		return nil, errVal.(error)
	} else {
		return goValueToPragmash(res[0].Interface())
	}
}

func valueToFloat(v *Value) (reflect.Value, error) {
	num, err := v.Number()
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(num.Float()), nil
}

func valueToFloatArray(v *Value) (reflect.Value, error) {
	valArr := v.Array()
	floats := make([]float64, len(valArr))
	for i, x := range valArr {
		num, err := x.Number()
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		floats[i] = num.Float()
	}
	return reflect.ValueOf(floats), nil
}

func valueToInt(v *Value) (reflect.Value, error) {
	num, err := v.Number()
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(int(num.Float())), nil
}

func valueToIntArray(v *Value) (reflect.Value, error) {
	valArr := v.Array()
	ints := make([]int, len(valArr))
	for i, x := range valArr {
		num, err := x.Number()
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		ints[i] = int(num.Float())
	}
	return reflect.ValueOf(ints), nil
}

func valueToNum(v *Value) (reflect.Value, error) {
	num, err := v.Number()
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(num), nil
}

func valueToNumArray(v *Value) (reflect.Value, error) {
	valArr := v.Array()
	nums := make([]*Number, len(valArr))
	for i, x := range valArr {
		num, err := x.Number()
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		nums[i] = num
	}
	return reflect.ValueOf(nums), nil
}

func valueToStrArray(v *Value) (reflect.Value, error) {
	valArr := v.Array()
	strs := make([]string, len(valArr))
	for i, x := range valArr {
		strs[i] = x.String()
	}
	return reflect.ValueOf(strs), nil
}
