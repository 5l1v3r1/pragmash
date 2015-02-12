package pragmash

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// StdGenerators implements ways of generating data.
type StdGenerators struct{}

// Echo joins its arguments with spaces.
func (s StdGenerators) Echo(args []Value) Value {
	interfaceArgs := make([]interface{}, len(args))
	for i, x := range args {
		interfaceArgs[i] = x
	}
	return StringValue(fmt.Sprint(interfaceArgs...))
}

// Range generates a range of integers.
func (s StdGenerators) Range(args []Value) (Value, error) {
	// Validate argument count.
	if len(args) == 0 || len(args) > 3 {
		return nil, errors.New("range cannot take " + strconv.Itoa(len(args)) +
			" arguments")
	}

	// Parse arguments.
	parsed := make([]int, len(args))
	for i, x := range args {
		var err error
		parsed[i], err = strconv.Atoi(x.String())
		if err != nil {
			return nil, err
		}
	}

	// Run the range function that corresponds to the number of arguments.
	if len(parsed) == 1 {
		return StringValue(rangeSingle(parsed[0])), nil
	} else if len(parsed) == 2 {
		return StringValue(rangeDouble(parsed[0], parsed[1])), nil
	} else {
		res, err := rangeTriple(parsed[0], parsed[1], parsed[2])
		if err != nil {
			return nil, err
		}
		return StringValue(res), nil
	}
}

func rangeDouble(start, end int) string {
	var buffer bytes.Buffer
	for i := start; i < end; i++ {
		if i != start {
			buffer.WriteRune('\n')
		}
		buffer.WriteString(strconv.Itoa(i))
	}
	return buffer.String()
}

func rangeSingle(end int) string {
	var buffer bytes.Buffer
	for i := 0; i < end; i++ {
		if i != 0 {
			buffer.WriteRune('\n')
		}
		buffer.WriteString(strconv.Itoa(i))
	}
	return buffer.String()
}

func rangeTriple(start, end, step int) (string, error) {
	if step == 0 {
		return "", errors.New("step cannot be 0")
	}

	var buffer bytes.Buffer
	i := start
	for {
		if step < 0 && start < end {
			break
		} else if step > 0 && start > end {
			break
		}
		if i != start {
			buffer.WriteRune('\n')
		}
		buffer.WriteString(strconv.Itoa(i))
		i += step
	}
	return buffer.String(), nil
}
