// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"github.com/go-pogo/errors"
	"reflect"
	"strconv"
)

func Marshal(v any) (Value, error) {
	return marshaler.Marshal(reflect.ValueOf(v))
}

type MarshalFunc func(v any) (string, error)

// GetMarshalFunc returns the globally registered MarshalFunc for reflect.Type
// typ or nil if there is none registered with RegisterMarshalFunc.
func GetMarshalFunc(typ reflect.Type) MarshalFunc { return marshaler.Func(typ) }

var marshaler Marshaler

type Marshaler struct {
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
	if !m.register.initialized() {
		// marshaler is always initialized
		return marshaler.Func(typ)
	}

	return m.register.find(typ)
}

// Marshal returns the string representation of the value.
// If the underlying reflect.Value is nil, it returns an empty string.
func (m *Marshaler) Marshal(val reflect.Value) (Value, error) {
	if fn := m.Func(val.Type()); fn != nil {
		return fn.Exec(val)
	}

	ot := val.Type()
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.String:
		return Value(val.String()), nil

	case reflect.Bool:
		return Value(strconv.FormatBool(val.Bool())), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Value(strconv.FormatInt(val.Int(), 10)), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Value(strconv.FormatUint(val.Uint(), 10)), nil

	case reflect.Float32, reflect.Float64:
		return Value(strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits())), nil

	case reflect.Complex64, reflect.Complex128:
		return Value(strconv.FormatComplex(val.Complex(), 'g', -1, val.Type().Bits())), nil
	}

	return "", errors.WithStack(&UnsupportedTypeError{Type: ot})
}

func (fn MarshalFunc) Exec(val reflect.Value) (Value, error) {
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
	}

	str, err := fn(val.Interface())
	if err != nil {
		return Value(str), errors.WithStack(err)
	}

	return Value(str), nil
}
