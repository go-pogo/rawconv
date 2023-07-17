// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"encoding"
	"fmt"
	"github.com/go-pogo/errors"
	"reflect"
	"strconv"
)

const (
	ParseError      errors.Kind = "parse error"
	ValidationError errors.Kind = "validation error"
)

func errKind(err error) errors.Kind {
	if ne, ok := err.(*strconv.NumError); ok {
		if ne.Err == strconv.ErrRange {
			return ValidationError
		} else {
			return ParseError
		}
	}
	return errors.UnknownKind
}

// Value is a textual representation of a value which is able to cast itself to
// any of the supported types using its corresponding method.
//
//	boolVal, err := parseval.Value("true").Bool()
type Value string

func From(v interface{}) (Value, error) {
	if rv, ok := v.(reflect.Value); ok {
		v = rv.Interface()
	}

	switch v := v.(type) {
	case encoding.TextMarshaler:
		b, err := v.MarshalText()
		return Value(b), errors.WithStack(err)
	case string:
		return Value(v), nil
	case []byte:
		return Value(v), nil
	}

	return Value(fmt.Sprintf("%v", v)), nil
}

// Empty indicates if Value is an empty string.
func (v Value) Empty() bool { return string(v) == "" }

func (v Value) GoString() string { return `parseval.Value("` + v.String() + `")` }

// String returns Value as a raw string.
func (v Value) String() string { return string(v) }

// StringVar sets the value p points to, to Value as raw string.
func (v Value) StringVar(p *string) { *p = v.String() }

// Bytes returns Value as raw bytes.
func (v Value) Bytes() []byte { return []byte(v) }

// BytesVar sets the value p points to, to Value as raw bytes.
func (v Value) BytesVar(p *[]byte) { *p = v.Bytes() }
