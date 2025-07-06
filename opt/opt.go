package opt

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Val is a generic type, which implements a field that can be one of three states:
//
// - field is not set in the request
// - field is explicitly set to `null` in the request
// - field is explicitly set to a valid value in the request
//
// Val is intended to be used with JSON marshalling and unmarshalling.
//
// Internal implementation details:
//
// - map[true]T means a value was provided
// - map[false]T means an explicit null was provided
// - nil or zero map means the field was not provided
//
// If the field is expected to be optional, add the `omitempty` JSON tags. Do NOT use `*Val`!
//
// Adapted from https://github.com/golang/go/issues/64515#issuecomment-1841057182
type Val[T any] map[bool]T

// New is a convenience helper to allow constructing a [Val] with a given value, for instance to construct a field inside a struct, without introducing an intermediate variable
func New[T any](t T) Val[T] {
	var n Val[T]
	n.Set(t)
	return n
}

// NewNull is a convenience helper to allow constructing a [Val] with an explicit `null`, for instance to construct a field inside a struct, without introducing an intermediate variable
func NewNull[T any]() Val[T] {
	var n Val[T]
	n.SetNull()
	return n
}

// Get retrieves the underlying value, if present, and returns an error if the value was not present
func (t Val[T]) Get() (T, error) {
	var empty T
	if t.IsNull() {
		return empty, errors.New("value is null")
	}
	if !t.IsSpecified() {
		return empty, errors.New("value is not specified")
	}
	return t[true], nil
}

// MustGet retrieves the underlying value, if present, and panics if the value was not present
func (t Val[T]) MustGet() T {
	v, err := t.Get()
	if err != nil {
		panic(err)
	}
	return v
}

// Set sets the underlying value to a given value
func (t *Val[T]) Set(value T) {
	*t = map[bool]T{true: value}
}

// IsNull indicate whether the field was sent, and had a value of `null`
func (t Val[T]) IsNull() bool {
	_, foundNull := t[false]
	return foundNull
}

// SetNull indicate that the field was sent, and had a value of `null`
func (t *Val[T]) SetNull() {
	var empty T
	*t = map[bool]T{false: empty}
}

// IsSpecified indicates whether the field was sent
func (t Val[T]) IsSpecified() bool {
	return len(t) != 0
}

// SetUnspecified indicate whether the field was sent
func (t *Val[T]) SetUnspecified() {
	*t = map[bool]T{}
}

func (t Val[T]) MarshalJSON() ([]byte, error) {
	// if field was specified, and `null`, marshal it
	if t.IsNull() {
		return []byte("null"), nil
	}

	// if field was unspecified, and `omitempty` is set on the field's tags, `json.Marshal` will omit this field

	// otherwise: we have a value, so marshal it
	return json.Marshal(t[true])
}

func (t *Val[T]) UnmarshalJSON(data []byte) error {
	// if field is unspecified, UnmarshalJSON won't be called

	// if field is specified, and `null`
	if bytes.Equal(data, []byte("null")) {
		t.SetNull()
		return nil
	}
	// otherwise, we have an actual value, so parse it
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.Set(v)
	return nil
}
