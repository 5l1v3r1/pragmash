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
func (_ StdArray) Arr(args ...*Value) []*Value {
	// Count the total number of elements so we only need one allocation.
	count := 0
	for _, v := range args {
		count += len(v.Array())
	}

	// Append the elements.
	values := make([]*Value, 0, count)
	for _, v := range args {
		values = append(values, v.Array()...)
	}
	return values
}

// Contains checks if an array contains a value.
func (_ StdArray) Contains(arr []string, val string) bool {
	for _, s := range arr {
		if s == val {
			return true
		}
	}
	return false
}

// Delete removes an element at a certain index from the array.
func (_ StdArray) Delete(arr []*Value, idx int) ([]*Value, error) {
	if idx < 0 || idx >= len(arr) {
		return nil, errors.New("index out of bounds: " + strconv.Itoa(idx))
	}
	res := make([]*Value, len(arr)-1)
	copy(res, arr[0:idx])
	copy(res[idx:], arr[idx+1:])
	return res, nil
}

// Insert inserts an element at a certain index in the array.
func (_ StdArray) Insert(arr []*Value, idx int, val *Value) ([]*Value, error) {
	if idx < 0 || idx > len(arr) {
		return nil, errors.New("index out of bounds: " + strconv.Itoa(idx))
	}
	res := make([]*Value, len(arr)+1)
	copy(res, arr[0:idx])
	copy(res[idx+1:], arr[idx:])
	res[idx] = val
	return res, nil
}

// Range generates a range of integers.
func (_ StdArray) Range(args ...int) ([]int, error) {
	// Validate argument count.
	if len(args) == 0 || len(args) > 3 {
		return nil, errors.New("range cannot take " + strconv.Itoa(len(args)) +
			" arguments")
	}

	// Run the range function that corresponds to the number of arguments.
	if len(args) == 1 {
		return rangeSingle(args[0]), nil
	} else if len(args) == 2 {
		return rangeDouble(args[0], args[1]), nil
	} else {
		return rangeTriple(args[0], args[1], args[2])
	}
}

// Shuffle randomly re-orders an array.
func (_ StdArray) Shuffle(list []*Value) []*Value {
	result := make([]*Value, len(list))
	perm := rand.Perm(len(list))
	for i, j := range perm {
		result[i] = list[j]
	}
	return result
}

// Sort sorts an array of strings alphabetically.
func (_ StdArray) Sort(arr []*Value) []*Value {
	cpy := make([]*Value, len(arr))
	copy(cpy, arr)
	sort.Sort(strValueList(cpy))
	return cpy
}

// Sortnums sorts an array of numbers.
func (_ StdArray) Sortnums(arr []*Value) ([]*Value, error) {
	// Make sure all the values are numbers.
	for _, v := range arr {
		if _, err := v.Number(); err != nil {
			return nil, err
		}
	}

	// Perform the sort.
	cpy := make([]*Value, len(arr))
	copy(cpy, arr)
	sort.Sort(numValueList(cpy))
	return cpy, nil
}

// Subarr returns a portion from an array.
func (_ StdArray) Subarr(arr []*Value, start, end int) []*Value {
	if len(arr) == 0 {
		return arr
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

	return arr[start:end]
}

// Sum takes arrays of numbers and returns their total sum.
func (_ StdArray) Sum(args ...[]*Number) *Number {
	sum := NewNumberInt(0)
	for _, list := range args {
		for _, num := range list {
			sum = AddNumbers(sum, num)
		}
	}
	return sum
}

func rangeDouble(start, end int) []int {
	res := make([]int, end-start)
	for i := start; i < end; i++ {
		res[i-start] = i
	}
	return res
}

func rangeSingle(end int) []int {
	res := make([]int, end)
	for i := 0; i < end; i++ {
		res[i] = i
	}
	return res
}

func rangeTriple(start, end, step int) ([]int, error) {
	if step == 0 {
		return nil, errors.New("step cannot be 0")
	}

	res := make([]int, 0)
	i := start
	for {
		if step < 0 && i <= end {
			break
		} else if step > 0 && i >= end {
			break
		}
		res = append(res, i)
		i += step
	}
	return res, nil
}

type numValueList []*Value

func (v numValueList) Len() int {
	return len(v)
}

func (v numValueList) Less(i, j int) bool {
	n1, _ := v[i].Number()
	n2, _ := v[j].Number()
	return CompareNumbers(n1, n2) < 0
}

func (v numValueList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

type strValueList []*Value

func (v strValueList) Len() int {
	return len(v)
}

func (v strValueList) Less(i, j int) bool {
	return v[i].String() < v[j].String()
}

func (v strValueList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
