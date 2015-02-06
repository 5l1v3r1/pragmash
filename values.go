package pragmash

type BoolValue bool

// A Value is a read-only variable value.
type Value interface {
    // Array returns the array representation of the value.
    Array() []Value

    // Context returns the context of the value. This is useful if the value is
    // an exception. In most cases, this should be an empty string.
    Context() string

    // Number returns the numerical representation of the value, or an error if
    // the value is not a number.
    Number() (Number, error)

    // String returns the textual representation of the value.
    String() string
}

