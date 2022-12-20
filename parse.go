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
	ErrUnableToAddr    errors.Msg = "unable to addr value"

	ParserError errors.Kind = "parser error"

	panicUnsupportedKind = "unsupported kind"
)

var dp = newParser(nil)

// Unmarshal Value val to any of the supported types.
func Unmarshal(val Value, v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New(ErrPointerExpected)
	}

	return dp.Parse(val, reflect.ValueOf(v))
}

// Parse Value v to any of the supported types and set its value to
// reflect.Value rval.
func Parse(val Value, dest reflect.Value) error { return dp.Parse(val, dest) }

type ParseFunc func(val Value, dest interface{}) error

func Register(typ reflect.Type, fn ParseFunc) { dp.Register(typ, fn) }

func init() {
	var (
		duration          = reflect.TypeOf(time.Nanosecond)
		urlType           = reflect.TypeOf(url.URL{})
		textUnmarshaler   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
		binaryUnmarshaler = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
	)

	// common types
	Register(duration, func(val Value, dest interface{}) (err error) {
		*dest.(*time.Duration), err = time.ParseDuration(val.String())
		return
	})
	Register(urlType, func(val Value, dest interface{}) error {
		u, err := url.ParseRequestURI(val.String())
		*dest.(*url.URL) = *u
		return err
	})

	// interfaces
	Register(textUnmarshaler, func(val Value, dest interface{}) error {
		return dest.(encoding.TextUnmarshaler).UnmarshalText(val.Bytes())
	})
	Register(binaryUnmarshaler, func(val Value, dest interface{}) error {
		return dest.(encoding.BinaryUnmarshaler).UnmarshalBinary(val.Bytes())
	})
}

type Parser struct {
	dp    *Parser
	types map[reflect.Kind]map[reflect.Type]int
	funcs []ParseFunc
}

func NewParser() *Parser { return newParser(dp) }

func newParser(dp *Parser) *Parser {
	return &Parser{
		dp:    dp,
		types: make(map[reflect.Kind]map[reflect.Type]int, 4),
		funcs: make([]ParseFunc, 0, 2),
	}
}

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

	if _, ok := p.types[k]; !ok {
		p.types[k] = map[reflect.Type]int{typ: len(p.funcs)}
	} else {
		p.types[k][typ] = len(p.funcs)
	}

	p.funcs = append(p.funcs, fn)
	return p
}

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "type `" + e.Type.String() + "` is not supported"
}

func (p *Parser) Func(typ reflect.Type) (ParseFunc, bool) {
	if kind, ok := p.types[typ.Kind()]; ok {
		if i, ok := kind[typ]; ok {
			return p.mustFunc(i), ok
		}
	}
	if p.dp != nil {
		return p.dp.Func(typ)
	}

	return nil, false
}

func (p *Parser) mustFunc(i int) ParseFunc {
	if i >= len(p.funcs) {
		panic("func should exist!")
	}
	return p.funcs[i]
}

// Parse Value v and set it to dest.
func (p *Parser) Parse(v Value, dest reflect.Value) error {
	// try exact type match
	if parseFn, ok := p.Func(dest.Type()); ok {
		return p.parse(v, dest, parseFn)
	}

	rv := dest
	for rv.Kind() == reflect.Ptr {
		// create a pointer to the type rval points to
		ptr := reflect.New(rv.Type().Elem())
		rv.Set(ptr)

		// take the value where the pointer points to and try parsing again...
		rv = ptr.Elem()
		if parseFn, ok := p.Func(rv.Type()); ok {
			return p.parse(v, rv, parseFn)
		}
	}

	// try interface implementations
	if rv, err := addr(rv); err != nil {
		return err
	} else {
		rt := rv.Type()
		for typ, i := range p.types[reflect.Interface] {
			if !rt.Implements(typ) {
				continue
			}

			if parseErr := p.parse(v, rv, p.mustFunc(i)); parseErr != nil {
				errors.Append(&err, parseErr)
			} else {
				return nil
			}
		}
	}

	// handle aliases of primitive types
	if v.Empty() {
		return nil
	}

	switch rv.Kind() {
	case reflect.String:
		rv.SetString(v.String())
		return nil

	case reflect.Bool:
		x, err := v.Bool()
		rv.SetBool(x)
		return err

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := intSize(v, rv.Type().Bits())
		rv.SetInt(x)
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := uintSize(v, rv.Type().Bits())
		rv.SetUint(x)
		return err

	case reflect.Float32, reflect.Float64:
		x, err := floatSize(v, rv.Type().Bits())
		rv.SetFloat(x)
		return err

	case reflect.Complex64, reflect.Complex128:
		x, err := complexSize(v, rv.Type().Bits())
		rv.SetComplex(x)
		return err
	}

	return errors.WithStack(&UnsupportedTypeError{Type: dest.Type()})
}

func (p *Parser) parse(v Value, rv reflect.Value, parseFn ParseFunc) error {
	rv, err := addr(rv)
	if err != nil {
		return errors.WithKind(err, ParserError)
	}

	dest := rv.Interface()
	if err = parseFn(v, dest); err != nil {
		if errors.GetKind(err) == errors.UnknownKind {
			return errors.WithKind(err, ParseError)
		} else {
			return err
		}
	}
	return nil
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
