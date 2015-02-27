package pragmash

import (
	"bytes"
	"strings"
)

var emptyValue = NewValueBool(false)

// A Value caches its various representations for use in programs.
// Accesses to a value's fields should be synchronized externally.
type Value struct {
	arrayRep  []*Value
	boolRep   bool
	numRep    *Number
	numErr    error
	stringRep *string
}

// NewValueArray creates a new Value from an array.
func NewValueArray(arr []*Value) *Value {
	// If there is one empty element, the array must be empty in order to
	// maintain integrity.
	if len(arr) == 1 && len(arr[0].String()) == 0 {
		str := ""
		return &Value{[]*Value{}, false, nil, nil, &str}
	}
	return &Value{arr, len(arr) != 0, nil, nil, nil}
}

// NewValueBool creates a new Value from a boolean.
func NewValueBool(b bool) *Value {
	res := &Value{nil, b, nil, nil, nil}
	if b {
		str := "true"
		res.stringRep = &str
	} else {
		str := ""
		res.stringRep = &str
	}
	res.arrayRep = []*Value{res}
	return res
}

// NewValueString creates a new HybridValue from a string.
func NewValueString(str string) *Value {
	return &Value{nil, len(str) > 0, nil, nil, &str}
}

// NewValueNumber creates a new Value from a *Number.
func NewValueNumber(num *Number) *Value {
	res := &Value{nil, true, num, nil, nil}
	res.arrayRep = []*Value{res}
	return res
}

// Array returns an array which represents the value.
func (h *Value) Array() []*Value {
	if h.arrayRep != nil {
		return h.arrayRep
	}

	// Generate an array by splitting the string into parts.
	strVal := h.String()
	if len(strVal) == 0 {
		h.arrayRep = []*Value{}
		return h.arrayRep
	}
	comps := strings.Split(strVal, "\n")
	res := make([]*Value, len(comps))
	for i, x := range comps {
		res[i] = NewValueString(x)
	}
	h.arrayRep = res
	return res
}

// Bool returns the pre-cached boolean representation of the value.
func (h *Value) Bool() bool {
	return h.boolRep
}

// Number returns the numerical representation of the value, parsing it as
// needed.
func (h *Value) Number() (*Number, error) {
	if h.numRep != nil || h.numErr != nil {
		return h.numRep, h.numErr
	}
	h.numRep, h.numErr = ParseNumber(h.String())
	return h.numRep, h.numErr
}

// Run returns v, nil.
func (v *Value) Run(r Runner) (*Value, *Breakout) {
    return v, nil
}

// String returns the string representation of the value.
func (h *Value) String() string {
	if h.stringRep != nil {
		return *h.stringRep
	} else if h.numRep != nil {
		str := h.numRep.String()
		h.stringRep = &str
		return str
	} else if h.arrayRep != nil {
		var buffer bytes.Buffer
		for i, v := range h.arrayRep {
			if i != 0 {
				buffer.WriteRune('\n')
			}
			buffer.WriteString(v.String())
		}
		str := buffer.String()
		h.stringRep = &str
		return str
	}
	panic("no way to generate a string representation")
	return ""
}
