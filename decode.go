// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"encoding"
	"net/url"
	"reflect"
	"time"

	"github.com/go-pogo/errors"
)

const (
	ErrPointerExpected errors.Msg = "expected a non-nil pointer to a value"
	ErrUnableToSet     errors.Msg = "unable to set value"
	ErrUnableToAddr    errors.Msg = "unable to addr value"

	InvalidActionError errors.Kind = "invalid action"
)

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Is(err error) bool {
	t, ok := err.(*UnsupportedTypeError)
	return ok && e.Type == t.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "type `" + e.Type.String() + "` is not supported"
}

// Unmarshal Value val to any of the supported types.
func Unmarshal(val Value, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.WithKind(ErrPointerExpected, InvalidActionError)
	}

	return unmarshaler.unmarshal(val, rv)
}

// UnmarshalReflect tries to unmarshal Value v to a supported type which matches
// dest, and sets the parsed value to it.
func UnmarshalReflect(v Value, dest reflect.Value) error {
	return unmarshaler.Unmarshal(v, dest)
}

type UnmarshalFunc func(val Value, dest interface{}) error

var unmarshaler = &Unmarshaler{root: true}

type Unmarshaler struct {
	root  bool
	types map[reflect.Kind]map[reflect.Type]int
	funcs []UnmarshalFunc
}

func unmarshalText(val Value, dest interface{}) error {
	return dest.(encoding.TextUnmarshaler).UnmarshalText(val.Bytes())
}

func init() {
	// interfaces
	unmarshaler.Register(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem(), unmarshalText)

	// common types
	unmarshaler.Register(reflect.TypeOf(time.Nanosecond), unmarshalDuration)
	unmarshaler.Register(reflect.TypeOf(url.URL{}), unmarshalUrl)
}

// Register ParseFunc for typ, making it available for Unmarshal, Parse and
// any Parser.
func Register(typ reflect.Type, fn UnmarshalFunc) { unmarshaler.Register(typ, fn) }

const panicUnsupportedKind = "unsupported kind"

// Register UnmarshalFunc for typ but only for this Unmarshaler.
func (u *Unmarshaler) Register(typ reflect.Type, fn UnmarshalFunc) *Unmarshaler {
	k := typ.Kind()
	if k == reflect.Invalid ||
		k == reflect.Uintptr ||
		k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.UnsafePointer ||
		// not yet supported
		k == reflect.Array || k == reflect.Map || k == reflect.Slice {
		panic(panicUnsupportedKind)
	}

	if u.types == nil {
		u.types = make(map[reflect.Kind]map[reflect.Type]int, 3)
	}
	if u.funcs == nil {
		u.funcs = make([]UnmarshalFunc, 0, 2)
	}

	if _, ok := u.types[k]; !ok {
		u.types[k] = map[reflect.Type]int{typ: len(u.funcs)}
	} else {
		u.types[k][typ] = len(u.funcs)
	}

	u.funcs = append(u.funcs, fn)
	return u
}

// Func returns the globally registered UnmarshalFunc for reflect.Type typ or
// nil if there is none registered with Register.
func Func(typ reflect.Type) UnmarshalFunc { return unmarshaler.Func(typ) }

// Func returns the (globally) registered UnmarshalFunc for reflect.Type typ or nil
// if there is none registered with RegisterType.
func (u *Unmarshaler) Func(typ reflect.Type) UnmarshalFunc {
	if u.funcs == nil {
		if u.root {
			return nil
		}

		// no custom ParseFunc registered for this Parser,
		// act like we're defaultParser...
		return unmarshaler.Func(typ)
	}

	fn := u.getFunc(typ)
	if fn == nil && typ.Kind() != reflect.Ptr {
		fn = u.getFuncFromImpl(reflect.New(typ).Type())
	}
	return fn
}

func (u *Unmarshaler) getFunc(typ reflect.Type) UnmarshalFunc {
	if fn := u.getFuncFromType(typ); fn != nil {
		return fn
	}

	if typ.Kind() == reflect.Ptr {
		if fn := u.getFunc(typ.Elem()); fn != nil {
			return fn
		}
		if fn := u.getFuncFromImpl(typ); fn != nil {
			return fn
		}
	}

	return nil
}

func (u *Unmarshaler) getFuncFromType(typ reflect.Type) UnmarshalFunc {
	if kind, ok := u.types[typ.Kind()]; ok {
		if i, ok := kind[typ]; ok {
			return u.getFuncFromIndex(i)
		}
	}
	if !u.root {
		return unmarshaler.getFuncFromType(typ)
	}

	return nil
}

func (u *Unmarshaler) getFuncFromImpl(typ reflect.Type) UnmarshalFunc {
	for x, i := range u.types[reflect.Interface] {
		if typ.Implements(x) {
			return u.getFuncFromIndex(i)
		}
	}
	if !u.root {
		return unmarshaler.getFuncFromImpl(typ)
	}

	return nil
}

const panicInvalidFuncIndex = "parseval.Parser: invalid index, func must exist!"

func (u *Unmarshaler) getFuncFromIndex(i int) UnmarshalFunc {
	if i >= len(u.funcs) {
		panic(panicInvalidFuncIndex)
	}
	return u.funcs[i]
}

// Unmarshal tries to unmarshal Value v to a supported type which matches dest,
// and sets the parsed value to it.
func (u *Unmarshaler) Unmarshal(v Value, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr && !dest.CanSet() {
		return errors.New(ErrUnableToSet)
	}
	return u.unmarshal(v, dest)
}

func (u *Unmarshaler) unmarshal(v Value, dest reflect.Value) error {
	if unmarshalFn := u.Func(dest.Type()); unmarshalFn != nil {
		return unmarshalFn.Exec(v, dest)
	}

	if v.Empty() {
		return nil
	}

	ot := dest.Type()

	var err error
	for dest.Kind() == reflect.Ptr {
		// take the value dest points to and make sure it is not nil
		if dest, err = value(dest.Elem()); err != nil {
			return err
		}
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
	}

	return errors.WithStack(&UnsupportedTypeError{Type: ot})
}

// Exec executes the UnmarshalFunc by taking the address of dest, and passing it
// as an interface to UnmarshalFunc. It will return an error when the address of
// reflect.Value dest cannot be taken, or when it is unable to set.
// Any error returned by UnmarshalFunc is wrapped with ParseError.
func (fn UnmarshalFunc) Exec(v Value, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		if !dest.CanAddr() {
			return errors.WithKind(ErrUnableToAddr, InvalidActionError)
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
		if errors.GetKind(err) == errors.UnknownKind {
			return errors.WithKind(err, ParseError)
		}
		return err
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
