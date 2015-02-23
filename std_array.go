package pragmash

import (
	"errors"
	"math/rand"
	"sort"
	"strconv"
)

// StdArray implements ways of manipulating or creating arrays
type StdArray struct{}

// Arr creates an array by combining arrays passed to it as arguments.
func (_ StdArray) Arr(args []Value) Value {
	values := make([]Value, 0)
	for _, v := range args {
		values = append(values, v.Array()...)
	}
	return NewHybridValueArray(values)
}

// Delete removes an element at a certain index from the array.
func (_ StdArray) Delete(arr []Value, idx int) (Value, error) {
	if idx < 0 || idx >= len(arr) {
		return nil, errors.New("index out of bounds: " + strconv.Itoa(idx))
	}
	res := make([]Value, len(arr)-1)
	copy(res, arr[0:idx])
	copy(res[idx:], arr[idx+1:])
	return NewHybridValueArray(res), nil
}

// Insert inserts an element at a certain index in the array.
func (_ StdArray) Insert(arr []Value, idx int, val Value) (Value, error) {
	if idx < 0 || idx > len(arr) {
		return nil, errors.New("index out of bounds: " + strconv.Itoa(idx))
	}
	res := make([]Value, len(arr)+1)
	copy(res, arr[0:idx])
	copy(res[idx+1:], arr[idx:])
	res[idx] = val
	return NewHybridValueArray(res), nil
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
		return NewHybridValueArray(rangeSingle(parsed[0])), nil
	} else if len(parsed) == 2 {
		return NewHybridValueArray(rangeDouble(parsed[0], parsed[1])), nil
	} else {
		res, err := rangeTriple(parsed[0], parsed[1], parsed[2])
		if err != nil {
			return nil, err
		}
		return NewHybridValueArray(res), nil
	}
}

// Shuffle randomly re-orders an array.
func (_ StdArray) Shuffle(arguments []Value) (Value, error) {
	if len(arguments) != 1 {
		return nil, errors.New("expected 1 argument")
	}
	list := arguments[0].Array()
	result := make([]Value, len(list))
	perm := rand.Perm(len(list))
	for i, j := range perm {
		result[i] = list[j]
	}
	return NewHybridValueArray(result), nil
}

// Sort sorts an array of strings alphabetically.
func (_ StdArray) Sort(arr []string) Value {
	// TODO: presereve the Values to keep cached representations
	
	cpy := make([]string, len(arr))
	copy(cpy, arr)
	sort.Strings(cpy)
	
	valArray := make([]Value, len(arr))
	for i, x := range cpy {
		valArray[i] = NewHybridValueString(x)
	}
	return NewHybridValueArray(valArray)
}

// Sortnums sorts an array of numbers.
func (_ StdArray) Sortnums(v Value) (Value, error) {
	valList := v.Array()
	numList := make(numberList, len(valList))
	for i, x := range valList {
		var err error
		numList[i], err = x.Number()
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(numList)
	
	valArray := make([]Value, len(numList))
	for i, x := range numList {
		valArray[i] = NewHybridValueNumber(x)
	}
	return NewHybridValueArray(valArray), nil
}

// Subarr returns a portion from an array.
func (_ StdArray) Subarr(arr []Value, start, end int) Value {
	if len(arr) == 0 {
		return emptyValue
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
	return NewHybridValueArray(res)
}

func rangeDouble(start, end int) []Value {
	res := make([]Value, end-start)
	for i := start; i < end; i++ {
		res[i-start] = NewNumberInt(int64(i))
	}
	return res
}

func rangeSingle(end int) []Value {
	res := make([]Value, end)
	for i := 0; i < end; i++ {
		res[i] = NewNumberInt(int64(i))
	}
	return res
}

func rangeTriple(start, end, step int) ([]Value, error) {
	if step == 0 {
		return nil, errors.New("step cannot be 0")
	}

	res := make([]Value, 0)
	i := start
	for {
		if step < 0 && i <= end {
			break
		} else if step > 0 && i >= end {
			break
		}
		res = append(res, NewNumberInt(int64(i)))
		i += step
	}
	return res, nil
}

type numberList []Number

func (n numberList) Len() int {
	return len(n)
}

func (n numberList) Less(i, j int) bool {
	return CompareNumbers(n[i], n[j]) < 0
}

func (n numberList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

