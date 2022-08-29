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

func Unmarshal(val Value, i interface{}) error {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return errors.New(ErrPointerExpected)
	}
	return parser.Parse(val, reflect.ValueOf(i))
}

func Parse(val Value, dest reflect.Value) error {
	return parser.Parse(val, dest)
}

type Parser struct {
	typ   reflect.Type
	parse func(v Value, u interface{}) error
}

func NewParser(typ reflect.Type, fn func(v Value, u interface{}) error) *Parser {
	return &Parser{
		typ:   typ,
		parse: fn,
	}
}

func (p *Parser) Parse(val Value, dest reflect.Value) error {
	typ := dest.Type()
	if p.typ != nil && typ == p.typ {
		return p.parse(val, dest.Interface())
	}

	return val.reflect(typ, dest)
}

var (
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	timeDurationType    = reflect.TypeOf(time.Nanosecond)
	urlUrlType          = reflect.TypeOf(url.URL{})
)

func (v Value) reflect(typ reflect.Type, dest reflect.Value) error {
	if typ.Implements(textUnmarshalerType) {
		return v.UnmarshalText(dest.Interface().(encoding.TextUnmarshaler))
	}
	if v.Empty() {
		return nil
	}

	switch typ {
	case timeDurationType:
		if x, err := v.Duration(); err != nil {
			return err
		} else {
			dest.Set(reflect.ValueOf(x))
			return nil
		}

	case urlUrlType:
		if x, err := v.Url(); err != nil {
			return err
		} else {
			dest.Set(reflect.ValueOf(x))
			return nil
		}
	}

	switch dest.Kind() {
	case reflect.String:
		dest.SetString(v.String())
		return nil

	case reflect.Bool:
		x, err := v.Bool()
		dest.SetBool(x)
		return err

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := intSize(v, typ.Bits())
		dest.SetInt(x)
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := uintSize(v, typ.Bits())
		dest.SetUint(x)
		return err

	case reflect.Float32, reflect.Float64:
		x, err := floatSize(v, typ.Bits())
		dest.SetFloat(x)
		return err

	case reflect.Complex64, reflect.Complex128:
		x, err := complexSize(v, typ.Bits())
		dest.SetComplex(x)
		return err

	default:
		return errors.WithStack(&UnsupportedError{Type: typ})
	}
}

type UnsupportedError struct {
	Type reflect.Type
}

func (e *UnsupportedError) Error() string {
	return "type `" + e.Type.String() + "` is not supported"
}
