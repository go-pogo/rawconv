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
	ErrPointerExpected errors.Msg = "expected a pointer to a value"
	ErrUnableToSet     errors.Msg = "unable to set value"
	ErrUnableToAddr    errors.Msg = "unable to addr value"

	panicUnsupportedKind = "unsupported kind"
)

var defaultParser = new(Parser)

// Unmarshal Value val to any of the supported types.
func Unmarshal(val Value, v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New(ErrPointerExpected)
	}

	return defaultParser.Parse(val, reflect.ValueOf(v))
}

type ParseFunc func(val Value, dest interface{}) error

// Exec executes the ParseFunc by taken the address of dest, and passing it as
// an interface to ParseFunc. It will return an error when the address of
// reflect.Value dest cannot be taken, or when it is unable to set.
// Any error returned by ParseFunc is wrapped with ParseError.
func (fn ParseFunc) Exec(val Value, dest reflect.Value) error {
	dest, err := addr(dest)
	if err != nil {
		return err
	}
	if !dest.Elem().CanSet() {
		return errors.New(ErrUnableToSet)
	}

	v := dest.Interface()
	if err = fn(val, v); err != nil {
		if errors.GetKind(err) == errors.UnknownKind {
			return errors.WithKind(err, ParseError)
		} else {
			return err
		}
	}
	return nil
}

type Parser struct {
	parent *Parser
	types  map[reflect.Kind]map[reflect.Type]int
	funcs  []ParseFunc
}

func NewParser() *Parser {
	return &Parser{parent: defaultParser}
}

func Register(typ reflect.Type, fn ParseFunc) { defaultParser.Register(typ, fn) }

func (p *Parser) Register(typ reflect.Type, fn ParseFunc) *Parser {
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

	if p.types == nil {
		p.types = make(map[reflect.Kind]map[reflect.Type]int, 3)
	}
	if p.funcs == nil {
		p.funcs = make([]ParseFunc, 0, 2)
	}

	if _, ok := p.types[k]; !ok {
		p.types[k] = map[reflect.Type]int{typ: len(p.funcs)}
	} else {
		p.types[k][typ] = len(p.funcs)
	}

	p.funcs = append(p.funcs, fn)
	return p
}

// HasFunc indicates if Parser has a ParseFunc registered for typ.
func (p *Parser) HasFunc(typ reflect.Type) bool {
	_, ok := p.Func(typ)
	return ok
}

// Func returns the registered ParseFunc for reflect.Type typ and true if it
// exists. Otherwise, it will return nil and false.
func (p *Parser) Func(typ reflect.Type) (ParseFunc, bool) {
	if kind, ok := p.types[typ.Kind()]; ok {
		if i, ok := kind[typ]; ok {
			return p.mustFunc(i), ok
		}
	}
	if p.parent != nil {
		return p.parent.Func(typ)
	}

	return nil, false
}

const panicParseFunc = "func should exist!"

func (p *Parser) mustFunc(i int) ParseFunc {
	if i >= len(p.funcs) {
		panic(panicParseFunc)
	}
	return p.funcs[i]
}

// Parse Value v to any of the supported types and set its value to
// reflect.Value dest.
func Parse(v Value, dest reflect.Value) error { return defaultParser.Parse(v, dest) }

// Parse Value v and set it to dest.
func (p *Parser) Parse(v Value, dest reflect.Value) error {
	// try exact type match
	if parseFn, ok := p.Func(dest.Type()); ok {
		return parseFn.Exec(v, dest)
	}

	ot := dest.Type()
	for dest.Kind() == reflect.Ptr {
		if dest.IsNil() {
			return errors.New(ErrPointerExpected)
		}

		// take the value where the pointer points to and try parsing again...
		dest = dest.Elem()
		if parseFn, ok := p.Func(dest.Type()); ok {
			return parseFn.Exec(v, dest)
		}
	}

	// try interface implementations
	if rv, err := addr(dest); err != nil {
		return err
	} else {
		rt := rv.Type()
		for typ, i := range p.types[reflect.Interface] {
			if !rt.Implements(typ) {
				continue
			}

			if parseErr := p.mustFunc(i).Exec(v, rv); parseErr != nil {
				errors.Append(&err, parseErr)
			} else {
				return nil
			}
		}
	}

	if v.Empty() {
		return nil
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

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "type `" + e.Type.String() + "` is not supported"
}

func addr(rv reflect.Value) (reflect.Value, error) {
	if rv.Kind() == reflect.Ptr {
		return rv, nil
	}
	if !rv.CanAddr() {
		return rv, errors.New(ErrUnableToAddr)
	}
	return rv.Addr(), nil
}

func init() {
	// interfaces
	Register(
		reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem(),
		func(val Value, dest interface{}) error {
			return dest.(encoding.TextUnmarshaler).UnmarshalText(val.Bytes())
		},
	)
	Register(
		reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem(),
		func(val Value, dest interface{}) error {
			return dest.(encoding.BinaryUnmarshaler).UnmarshalBinary(val.Bytes())
		},
	)

	// common types
	Register(
		reflect.TypeOf(time.Nanosecond),
		func(val Value, dest interface{}) (err error) {
			*dest.(*time.Duration), err = time.ParseDuration(val.String())
			return
		},
	)
	Register(
		reflect.TypeOf(url.URL{}),
		func(val Value, dest interface{}) error {
			u, err := url.ParseRequestURI(val.String())
			*dest.(*url.URL) = *u
			return err
		},
	)
}
