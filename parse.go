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

var defaultParser = new(Parser)

func init() {
	// interfaces
	Register(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem(), unmarshalText)
	Register(reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem(), unmarshalBinary)

	// common types
	Register(reflect.TypeOf(time.Nanosecond), parseDuration)
	Register(reflect.TypeOf(url.URL{}), parseUrl)
}

// Unmarshal Value val to any of the supported types.
func Unmarshal(val Value, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.WithKind(ErrPointerExpected, InvalidActionError)
	}

	return defaultParser.parse(val, rv)
}

type ParseFunc func(val Value, dest interface{}) error

func unmarshalText(val Value, dest interface{}) error {
	return dest.(encoding.TextUnmarshaler).UnmarshalText(val.Bytes())
}

func unmarshalBinary(val Value, dest interface{}) error {
	return dest.(encoding.BinaryUnmarshaler).UnmarshalBinary(val.Bytes())
}

func parseDuration(val Value, dest interface{}) (err error) {
	*dest.(*time.Duration), err = time.ParseDuration(val.String())
	return
}

func parseUrl(val Value, dest interface{}) error {
	u, err := url.ParseRequestURI(val.String())
	*dest.(*url.URL) = *u
	return err
}

// Exec executes the ParseFunc by taken the address of dest, and passing it as
// an interface to ParseFunc. It will return an error when the address of
// reflect.Value dest cannot be taken, or when it is unable to set.
// Any error returned by ParseFunc is wrapped with ParseError.
func (fn ParseFunc) Exec(val Value, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		if !dest.CanAddr() {
			return errors.WithKind(ErrUnableToAddr, InvalidActionError)
		}
		return fn.exec(val, dest.Addr())
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

	return fn.exec(val, dest)
}

func (fn ParseFunc) exec(val Value, dest reflect.Value) error {
	if err := fn(val, dest.Interface()); err != nil {
		if errors.GetKind(err) == errors.UnknownKind {
			return errors.WithKind(err, ParseError)
		}
		return err
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

// Register ParseFunc for typ, making it available for Unmarshal, Parse and
// any Parser.
func Register(typ reflect.Type, fn ParseFunc) { defaultParser.Register(typ, fn) }

const panicUnsupportedKind = "unsupported kind"

// Register ParseFunc for typ but only for this Parser.
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

// Func returns the ParseFunc for reflect.Type typ and true if it is globally
// registered with Register. Otherwise, it will return nil and false.
func Func(typ reflect.Type) (ParseFunc, bool) {
	return defaultParser.Func(typ)
}

// Func returns the ParseFunc for reflect.Type typ and true if it is (globally)
// registered with Register. Otherwise, it will return nil and false.
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

const panicInvalidFuncIndex = "parseval.Parser: invalid index, func must exist!"

func (p *Parser) mustFunc(i int) ParseFunc {
	if i >= len(p.funcs) {
		panic(panicInvalidFuncIndex)
	}
	return p.funcs[i]
}

// Parse Value v to any of the supported types and set its value to
// reflect.Value dest.
func Parse(v Value, dest reflect.Value) error { return defaultParser.Parse(v, dest) }

// Parse Value v to any of the supported types and set its value to
// reflect.Value dest.
func (p *Parser) Parse(v Value, dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr && !dest.CanSet() {
		return errors.New(ErrUnableToSet)
	}
	return p.parse(v, dest)
}

func (p *Parser) parse(v Value, dest reflect.Value) error {
	// try exact type match
	if parseFn, ok := p.Func(dest.Type()); ok {
		return parseFn.Exec(v, dest)
	}

	ot := dest.Type()

	var err error
	for dest.Kind() == reflect.Ptr {
		// take the value dest points to and make sure it is not nil
		if dest, err = value(dest.Elem()); err != nil {
			return err
		}

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

// addr returns a pointer to the reflect.Value.
func addr(rv reflect.Value) (reflect.Value, error) {
	if rv.Kind() == reflect.Ptr {
		return rv, nil
	}
	if !rv.CanAddr() {
		return rv, errors.WithKind(ErrUnableToAddr, InvalidActionError)
	}
	return rv.Addr(), nil
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
