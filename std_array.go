package pragmash

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// StdArray implements ways of manipulating or creating arrays
type StdArray struct{}

// Arr joins its arguments with newlines.
func (_ StdArray) Arr(args []Value) Value {
	var buffer bytes.Buffer
	for i, v := range args {
		if i != 0 {
			buffer.WriteRune('\n')
		}
		buffer.WriteString(v.String())
	}
	return StringValue(buffer.String())
}

// Delete removes an element at a certain index from the array.
func (_ StdArray) Delete(arr []string, idx int) (Value, error) {
	if idx < 0 || idx >= len(arr) {
		return nil, errors.New("index out of bounds: " + strconv.Itoa(idx))
	}
	res := make([]string, len(arr)-1)
	copy(res, arr[0:idx])
	copy(res[idx:], arr[idx+1:])
	return StringValue(strings.Join(res, "\n")), nil
}

// Insert inserts an element at a certain index in the array.
func (_ StdArray) Insert(arr []string, idx int, val string) (Value, error) {
	if idx < 0 || idx > len(arr) {
		return nil, errors.New("index out of bounds: " + strconv.Itoa(idx))
	}
	res := make([]string, len(arr)+1)
	copy(res, arr[0:idx])
	copy(res[idx+1:], arr[idx:])
	res[idx] = val
	return StringValue(strings.Join(res, "\n")), nil
}

// Range generates a range of integers.
func (_ StdArray) Range(args []Value) (Value, error) {
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

// Subarr returns a subarray.
func (_ StdArray) Subarr(arr []string, start, end int) Value {
	if len(arr) == 0 {
		return StringValue("")
	}
	
	// Sanitize the range
	if start < 0 {
		start = 0
	} else if start > len(arr) {
		start = len(arr)
	}
	if end < start {
		end = start
	} else if end > len(arr) {
		end = len(arr)
	}
	
	res := arr[start : end]
	return StringValue(strings.Join(res, "\n"))
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
		if step < 0 && i <= end {
			break
		} else if step > 0 && i >= end {
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
