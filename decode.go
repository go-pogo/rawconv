// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"reflect"
	"strings"

	"github.com/go-pogo/errors"
)

const (
	ErrPointerExpected    errors.Msg = "expected a non-nil pointer to a value"
	ErrUnmarshalNested    errors.Msg = "cannot unmarshal nested array/slice/map"
	ErrUnableToSet        errors.Msg = "unable to set value"
	ErrUnableToAddr       errors.Msg = "unable to addr value"
	ErrRuneTooManyChars   errors.Msg = "too many characters"
	ErrArrayTooManyValues errors.Msg = "too many values"
	ErrMapInvalidFormat   errors.Msg = "invalid map format"
	ErrUnmarshalFuncExec  errors.Msg = "error while executing UnmarshalFunc"
)

// Unmarshal parses Value and stores the result in the value pointed to by v.
// If v is nil or not a pointer, Unmarshal returns an ErrPointerExpected error.
// If v is not a supported type an UnsupportedTypeError is returned.
// By default, the following types are supported:
//   - string
//   - bool
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - complex64, complex128
//   - array, slice
//   - map
//   - time.Duration
//   - url.URL
//   - encoding.TextUnmarshaler
//
// Use RegisterUnmarshalFunc to add additional (custom) types.
func Unmarshal(val Value, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New(ErrPointerExpected)
	}

	return unmarshaler.unmarshal(val, rv, false)
}

// UnmarshalFunc is a function which can unmarshal a Value to any type.
// Argument dest is always a pointer to the value to unmarshal to.
type UnmarshalFunc func(val Value, dest any) error

// GetUnmarshalFunc returns the globally registered UnmarshalFunc for
// reflect.Type typ or nil if there is none registered with
// RegisterUnmarshalFunc.
func GetUnmarshalFunc(typ reflect.Type) UnmarshalFunc { return unmarshaler.Func(typ) }

// unmarshaler is the global Unmarshaler.
var unmarshaler Unmarshaler

// Unmarshaler is a type which can unmarshal a Value to any type that's
// registered with Register. It wil always fallback to the global Unmarshaler
// when a type is not registered.
type Unmarshaler struct {
	Options
	register register[UnmarshalFunc]
}

// Register the UnmarshalFunc for typ but only for this Unmarshaler.
func (u *Unmarshaler) Register(typ reflect.Type, fn UnmarshalFunc) *Unmarshaler {
	u.register.add(typ, fn)
	return u
}

// Func returns the (globally) registered UnmarshalFunc for reflect.Type typ or
// nil if there is none registered with Register or RegisterUnmarshalFunc.
func (u *Unmarshaler) Func(typ reflect.Type) UnmarshalFunc {
	if u.register.initialized() {
		if fn := u.register.find(typ); fn != nil {
			return fn
		}
	}
	// fallback to global unmarshaler
	return unmarshaler.register.find(typ)
}

// Unmarshal tries to unmarshal Value to a supported type which matches the
// type of v, and sets the parsed value to it. See Unmarshal for additional
// details.
func (u *Unmarshaler) Unmarshal(val Value, v reflect.Value) error {
	if v.Kind() != reflect.Ptr && !v.CanSet() {
		return errors.New(ErrUnableToSet)
	}
	return u.unmarshal(val, v, false)
}

func (u *Unmarshaler) unmarshal(v Value, dest reflect.Value, nested bool) error {
	if fn := u.Func(dest.Type()); fn != nil {
		return fn.Exec(v, dest)
	}

	if v.IsEmpty() {
		return nil
	}

	ot := dest.Type()

	var err error
	for dest.Kind() == reflect.Ptr {
		if dest.IsNil() {
			if !dest.CanSet() {
				return errors.New(ErrUnableToSet)
			}

			dest.Set(reflect.New(dest.Type().Elem()))
		}
		dest = dest.Elem()
	}

	// handle aliases of primitive types
	switch dest.Kind() {
	case reflect.String:
		dest.SetString(v.String())
		return nil

	case reflect.Bool:
		x, err := v.Bool()
		dest.SetBool(x)
		return err

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := intSize(v, dest.Type().Bits())
		dest.SetInt(x)
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := uintSize(v, dest.Type().Bits())
		dest.SetUint(x)
		return err

	case reflect.Float32, reflect.Float64:
		x, err := floatSize(v, dest.Type().Bits())
		dest.SetFloat(x)
		return err

	case reflect.Complex64, reflect.Complex128:
		x, err := complexSize(v, dest.Type().Bits())
		dest.SetComplex(x)
		return err

	case reflect.Array:
		if nested {
			return errors.New(ErrUnmarshalNested)
		}

		parts := split(v.String(), u.itemSeparator())
		typ := dest.Type().Elem()

		partsLen, arrayLen := len(parts), dest.Len()
		for i := 0; i < partsLen && i < arrayLen; i++ {
			part := strings.TrimSpace(parts[i])
			val := reflect.New(typ).Elem()
			if err = u.unmarshal(Value(part), val, true); err != nil {
				return err
			}
			dest.Index(i).Set(val)
		}
		if partsLen > arrayLen {
			return errors.New(ErrArrayTooManyValues)
		}
		return nil

	case reflect.Slice:
		if nested {
			return errors.New(ErrUnmarshalNested)
		}

		parts := split(v.String(), u.itemSeparator())
		slice := reflect.MakeSlice(dest.Type(), 0, len(parts))
		typ := dest.Type().Elem()

		for _, part := range parts {
			part = strings.TrimSpace(part)
			val := reflect.New(typ).Elem()
			if err = u.unmarshal(Value(part), val, true); err != nil {
				return err
			}
			slice = reflect.Append(slice, val)
		}

		dest.Set(slice)
		return nil

	case reflect.Map:
		if nested {
			return errors.New(ErrUnmarshalNested)
		}

		parts := split(v.String(), u.itemSeparator())
		if dest.IsNil() {
			dest.Set(reflect.MakeMapWithSize(dest.Type(), len(parts)))
		}

		keyTyp := dest.Type().Key()
		valTyp := dest.Type().Elem()

		for _, part := range parts {
			kv := strings.SplitN(part, u.keyValueSeparator(), 2)
			if len(kv) != 2 {
				return errors.New(ErrMapInvalidFormat)
			}

			key := reflect.New(keyTyp).Elem()
			if err = u.unmarshal(Value(kv[0]), key, true); err != nil {
				return err
			}
			val := reflect.New(valTyp).Elem()
			if err = u.unmarshal(Value(kv[1]), val, true); err != nil {
				return err
			}

			dest.SetMapIndex(key, val)
		}
		return nil

	default:
		return errors.WithStack(&UnsupportedTypeError{Type: ot})
	}
}

// Exec executes the UnmarshalFunc by taking the address of dest, and passing it
// as an interface to UnmarshalFunc. It will return an error when the address of
// reflect.Value dest cannot be taken, or when it is unable to set.
// Any error returned by UnmarshalFunc is wrapped with ErrParseFailure.
func (fn UnmarshalFunc) Exec(v Value, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		if !dest.CanAddr() {
			return errors.New(ErrUnableToAddr)
		}
		return fn.exec(v, dest.Addr())
	}

	var err error
	if dest, err = value(dest); err != nil {
		return err
	}
	for dest.Elem().Kind() == reflect.Ptr {
		if dest, err = value(dest.Elem()); err != nil {
			return err
		}
	}

	return fn.exec(v, dest)
}

func (fn UnmarshalFunc) exec(val Value, dest reflect.Value) error {
	if err := fn(val, dest.Interface()); err != nil {
		return errors.Wrap(err, ErrUnmarshalFuncExec)
	}
	return nil
}

// value ensures the reflect.Value is not nil.
func value(rv reflect.Value) (reflect.Value, error) {
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		if !rv.CanSet() {
			return rv, errors.New(ErrUnableToSet)
		}

		rv.Set(reflect.New(rv.Type().Elem()))
	}
	return rv, nil
}

func split(str, sep string) []string {
	return strings.Split(str, sep)
}
