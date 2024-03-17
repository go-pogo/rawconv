// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"github.com/go-pogo/errors"
	"reflect"
	"strconv"
	"strings"
)

const ErrMarshalNested errors.Msg = "cannot marshal nested array/slice/map"

// Marshal formats the value pointed to by v to a raw string Value.
// If v is not a supported type an UnsupportedTypeError is returned.
// By default, the following types are supported:
// - encoding.TextUnmarshaler
// - time.Duration
// - url.URL
// - string
// - bool
// - int, int8, int16, int32, int64
// - uint, uint8, uint16, uint32, uint64
// - float32, float64
// - complex64, complex128
// Use RegisterMarshalFunc to add additional (custom) types.
func Marshal(v any) (Value, error) {
	return marshaler.Marshal(reflect.ValueOf(v))
}

type MarshalFunc func(v any) (string, error)

// GetMarshalFunc returns the globally registered MarshalFunc for reflect.Type
// typ or nil if there is none registered with RegisterMarshalFunc.
func GetMarshalFunc(typ reflect.Type) MarshalFunc { return marshaler.Func(typ) }

// marshaler is the global Marshaler.
var marshaler Marshaler

// Marshaler is a type which can marshal any reflect.Value to its raw string
// representation as long as it's registered with Register. It wil always
// fallback to the global Marshaler when a type is not registered.
type Marshaler struct {
	Options
	register register[MarshalFunc]
}

// Register the MarshalFunc for typ but only for this Marshaler.
func (m *Marshaler) Register(typ reflect.Type, fn MarshalFunc) *Marshaler {
	m.register.add(typ, fn)
	return m
}

// Func returns the (globally) registered MarshalFunc for reflect.Type typ or
// nil if there is none registered with Register or RegisterMarshalFunc.
func (m *Marshaler) Func(typ reflect.Type) MarshalFunc {
	if m.register.initialized() {
		if fn := m.register.find(typ); fn != nil {
			return fn
		}
	}
	// fallback to global marshaler
	return marshaler.register.find(typ)
}

// Marshal returns the string representation of the value.
// If the underlying reflect.Value is nil, it returns an empty string.
func (m *Marshaler) Marshal(val reflect.Value) (Value, error) {
	str, err := m.marshal(val, false)
	return Value(str), err
}

func (m *Marshaler) marshal(val reflect.Value, nested bool) (string, error) {
	if fn := m.Func(val.Type()); fn != nil {
		return fn.exec(val)
	}

	ot := val.Type()
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
	}

	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch val.Kind() {
	case reflect.String:
		return val.String(), nil

	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, 64), nil

	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(val.Complex(), 'g', -1, 128), nil

	case reflect.Array, reflect.Slice:
		if nested {
			return "", errors.New(ErrMarshalNested)
		}

		sep := m.itemSeparator()

		var buf strings.Builder
		for i := 0; i < val.Len(); i++ {
			v, err := m.marshal(val.Index(i), true)
			if err != nil {
				return "", err
			}

			if i > 0 {
				buf.WriteString(sep)
			}
			buf.WriteString(v)
		}
		return buf.String(), nil

	case reflect.Map:
		if nested {
			return "", errors.New(ErrMarshalNested)
		}

		sep1 := m.keyValueSeparator()
		sep2 := m.itemSeparator()

		var buf strings.Builder
		var firstDone bool
		for iter := val.MapRange(); iter.Next(); {
			v, err := m.marshal(iter.Value(), true)
			if err != nil {
				return "", err
			}
			k, err := m.marshal(iter.Key(), true)
			if err != nil {
				return "", err
			}

			if firstDone {
				buf.WriteString(sep2)
			}

			buf.WriteString(k)
			buf.WriteString(sep1)
			buf.WriteString(v)
			firstDone = true
		}
		return buf.String(), nil
	}

	return "", errors.WithStack(&UnsupportedTypeError{Type: ot})
}

// Exec executes the MarshalFunc for the given reflect.Value.
func (fn MarshalFunc) Exec(val reflect.Value) (Value, error) {
	str, err := fn.exec(val)
	return Value(str), err
}

func (fn MarshalFunc) exec(val reflect.Value) (string, error) {
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
	}

	str, err := fn(val.Interface())
	if err != nil {
		return str, errors.WithStack(err)
	}

	return str, nil
}
