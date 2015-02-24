package pragmash

import (
	"bytes"
	"errors"
	"strings"
)

var emptyValue = BoolValue(false)

// A BoolValue is a bool which implements the Value interface.
type BoolValue bool

// Array returns an empty array if the receiver is false, or an array with a
// single true BoolValue if it's true.
func (b BoolValue) Array() []Value {
	if !b {
		return []Value{}
	} else {
		return []Value{b}
	}
}

// Bool returns bool(b).
func (b BoolValue) Bool() bool {
	return bool(b)
}

// Context returns an empty string.
func (b BoolValue) Context() string {
	return ""
}

// Number returns an error, since a boolean is not a number.
func (b BoolValue) Number() (Number, error) {
	return nil, errors.New("invalid number: " + b.String())
}

// String returns StringValue("true") for a true receiver and StringValue("")
// for a false one.
func (b BoolValue) String() string {
	if b {
		return "true"
	} else {
		return ""
	}
}

// A HybridValue is a Value which caches its various representations.
type HybridValue struct {
	ArrayRep  []Value
	BoolRep   bool
	NumVal    Number
	NumErr    error
	StringRep *string
}

// NewHybridValueArray creates a new HybridValue from an array of values.
func NewHybridValueArray(arr []Value) *HybridValue {
	// If there is one empty element, the array must be empty in order to
	// maintain integrity.
	if len(arr) == 1 && len(arr[0].String()) == 0 {
		str := ""
		return &HybridValue{[]Value{}, false, nil, nil, &str}
	}
	return &HybridValue{arr, len(arr) != 0, nil, nil, nil}
}

// NewHybridValueString creates a new HybridValue from a string.
func NewHybridValueString(str string) *HybridValue {
	return &HybridValue{nil, len(str) > 0, nil, nil, &str}
}

// NewHybridValueNumber creates a new HybridValue from a Number.
func NewHybridValueNumber(num Number) *HybridValue {
	res := &HybridValue{nil, true, num, nil, nil}
	res.ArrayRep = []Value{res}
	return res
}

// Array returns an array which represents the value.
func (h *HybridValue) Array() []Value {
	if h.ArrayRep != nil {
		return h.ArrayRep
	}
	
	// Generate an array by splitting the string into parts.
	strVal := h.String()
	if len(strVal) == 0 {
		h.ArrayRep = []Value{}
		return h.ArrayRep
	}
	comps := strings.Split(strVal, "\n")
	res := make([]Value, len(comps))
	for i, x := range comps {
		res[i] = NewHybridValueString(x)
	}
	h.ArrayRep = res
	return res
}

// Bool returns the pre-cached boolean representation of the value.
func (h *HybridValue) Bool() bool {
	return h.BoolRep
}

// Context returns an empty string.
func (_ *HybridValue) Context() string {
	return ""
}

// Number returns the numerical representation of the value, parsing it as
// needed.
func (h *HybridValue) Number() (Number, error) {
	if h.NumVal != nil || h.NumErr != nil {
		return h.NumVal, h.NumErr
	}
	h.NumVal, h.NumErr = ParseNumber(h.String())
	return h.NumVal, h.NumErr
}

// String returns the string representation of the value.
func (h *HybridValue) String() string {
	if h.StringRep != nil {
		return *h.StringRep
	} else if h.NumVal != nil {
		str := h.NumVal.String()
		h.StringRep = &str
		return str
	} else if h.ArrayRep != nil {
		var buffer bytes.Buffer
		for i, v := range h.ArrayRep {
			str := v.String()
			if i != 0 {
				buffer.WriteRune('\n')
			}
			buffer.WriteString(str)
		}
		str := buffer.String()
		h.StringRep = &str
		return str
	}
	panic("no way to generate a string representation")
	return ""
}

// A Value is a read-only variable value.
type Value interface {
	// Array returns the array representation of the value.
	Array() []Value

	// Bool returns the boolean representation of the value.
	Bool() bool

	// Context returns the context of the value. This is useful if the value is
	// an exception. In most cases, this should be an empty string.
	Context() string

	// Number returns the numerical representation of the value, or an error if
	// the value is not a number.
	Number() (Number, error)

	// String returns the textual representation of the value.
	String() string
}
