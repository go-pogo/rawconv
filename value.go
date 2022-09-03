// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"encoding"
	"net/url"
	"strconv"
	"time"

	"github.com/go-pogo/errors"
)

const (
	ParseError      errors.Kind = "parse error"
	ValidationError errors.Kind = "validation error"
)

var Separator string = ","

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

// Empty indicates if Value is an empty string.
func (v Value) Empty() bool { return v.String() == "" }

func (v Value) GoString() string { return `parseval.Value("` + v.String() + `")` }

// String returns Value as a raw string.
func (v Value) String() string { return string(v) }

// StringVar sets the value p points to, to Value as raw string.
func (v Value) StringVar(p *string) { *p = v.String() }

// Bytes returns Value as raw bytes.
func (v Value) Bytes() []byte { return []byte(v) }

// BytesVar sets the value p points to, to Value as raw bytes.
func (v Value) BytesVar(p *[]byte) { *p = v.Bytes() }

// Duration tries to parse Value as time.Duration with time.ParseDuration.
func (v Value) Duration() (time.Duration, error) {
	x, err := time.ParseDuration(v.String())
	return x, errors.WithKind(err, ParseError)
}

// DurationVar sets the value p points to using Duration.
func (v Value) DurationVar(p *time.Duration) (err error) {
	*p, err = v.Duration()
	return
}

// Url tries to parse Value as an url.Url with url.ParseRequestURI.
func (v Value) Url() (*url.URL, error) {
	x, err := url.ParseRequestURI(v.String())
	if err != nil {
		err = errors.WithKind(err, ParseError)
	}
	return x, err
}

// UrlVar sets the value p points to using Url.
func (v Value) UrlVar(p **url.URL) (err error) {
	*p, err = v.Url()
	return
}

// UnmarshalTextWith uses encoding.TextUnmarshaler u to unmarshal v.
func (v Value) UnmarshalTextWith(u encoding.TextUnmarshaler) error {
	err := u.UnmarshalText(v.Bytes())
	return errors.WithKind(err, errKind(err))
}

func (v Value) unmarshal(i interface{}) (bool, error) {
	if u, ok := i.(encoding.TextUnmarshaler); ok {
		// let TextUnmarshaler decide what to do with possible empty v
		return true, v.UnmarshalTextWith(u)
	}
	if v.Empty() {
		return true, nil
	}

	switch p := i.(type) {
	case *string:
		v.StringVar(p)
		return true, nil
	case *[]byte:
		v.BytesVar(p)
		return true, nil

	case *bool:
		return true, v.BoolVar(p)

	case *int:
		return true, v.IntVar(p)
	case *int8:
		return true, v.Int8Var(p)
	case *int16:
		return true, v.Int16Var(p)
	case *int32:
		return true, v.Int32Var(p)
	case *int64:
		return true, v.Int64Var(p)

	case *uint:
		return true, v.UintVar(p)
	case *uint8:
		return true, v.Uint8Var(p)
	case *uint16:
		return true, v.Uint16Var(p)
	case *uint32:
		return true, v.Uint32Var(p)
	case *uint64:
		return true, v.Uint64Var(p)

	case *float32:
		return true, v.Float32Var(p)
	case *float64:
		return true, v.Float64Var(p)

	case *complex64:
		return true, v.Complex64Var(p)
	case *complex128:
		return true, v.Complex128Var(p)

	case *time.Duration:
		return true, v.DurationVar(p)
	case *url.URL:
		return true, v.UrlVar(&p)
	}

	return false, nil
}
