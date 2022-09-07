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

const ErrPointerExpected errors.Msg = "expected a pointer to a value"

var parser = new(Parser)

// Unmarshal Value v to any of the supported types.
func Unmarshal(v Value, i interface{}) error {
	if ok, err := v.unmarshal(i); ok || err != nil {
		return err
	}
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return errors.New(ErrPointerExpected)
	}

	return parser.Parse(v, reflect.ValueOf(i))
}

// Parse Value v to any of the supported types and set its value to
// reflect.Value rval.
func Parse(v Value, rval reflect.Value) error { return parser.Parse(v, rval) }

type Parser struct {
	typ reflect.Type
	cp  func(Value, interface{}) error
}

func NewParser(rtyp reflect.Type, fn func(v Value, i interface{}) error) *Parser {
	return &Parser{
		typ: rtyp,
		cp:  fn,
	}
}

var (
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	timeDurationType    = reflect.TypeOf(time.Nanosecond)
	urlUrlType          = reflect.TypeOf(url.URL{})
)

// Parse Value v and set it to val.
func (p *Parser) Parse(v Value, rval reflect.Value) error {
	rtyp := rval.Type()
	if p.typ != nil && rtyp == p.typ {
		return p.cp(v, rval.Interface())
	}

	// for rval.Kind() == reflect.Ptr {
	// 	rval = rval.Elem()
	// }

	rtyp = rval.Type()
	if rtyp.Implements(textUnmarshalerType) {
		// let TextUnmarshaler decide what to do with possible empty v
		return v.UnmarshalTextWith(rval.Interface().(encoding.TextUnmarshaler))
	}
	if v.Empty() {
		return nil
	}

	// handle known types
	switch rtyp {
	case timeDurationType:
		if x, err := v.Duration(); err != nil {
			return err
		} else {
			rval.Set(reflect.ValueOf(x))
			return nil
		}

	case urlUrlType:
		if x, err := v.Url(); err != nil {
			return err
		} else {
			rval.Set(reflect.ValueOf(x))
			return nil
		}
	}

	// handle aliases of primitive types
	switch rval.Kind() {
	case reflect.String:
		rval.SetString(v.String())
		return nil

	case reflect.Bool:
		x, err := v.Bool()
		rval.SetBool(x)
		return err

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := intSize(v, rval.Type().Bits())
		rval.SetInt(x)
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := uintSize(v, rval.Type().Bits())
		rval.SetUint(x)
		return err

	case reflect.Float32, reflect.Float64:
		x, err := floatSize(v, rval.Type().Bits())
		rval.SetFloat(x)
		return err

	case reflect.Complex64, reflect.Complex128:
		x, err := complexSize(v, rval.Type().Bits())
		rval.SetComplex(x)
		return err
	}

	return errors.WithStack(&UnsupportedTypeError{Type: rtyp})
}

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "type `" + e.Type.String() + "` is not supported"
}
