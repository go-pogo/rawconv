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

// Value is a textual representation of a value which is able to cast itself to
// any of the supported types using its corresponding method.
//
//	boolVal, err := parseval.Value("true").Bool()
type Value string

// Empty indicates if Value is an empty string.
func (v Value) Empty() bool { return string(v) == "" }

func (v Value) GoString() string { return `parseval.Value("` + string(v) + `")` }

// String returns Value as a raw string.
func (v Value) String() string { return string(v) }

// Bool tries to parse Value as a bool with strconv.ParseBool.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns an error.
func (v Value) Bool() (bool, error) {
	x, err := strconv.ParseBool(string(v))
	return x, errors.WithKind(err, errKind(err))
}

// Int tries to parse Value as an int with strconv.ParseInt.
func (v Value) Int() (int, error) {
	x, err := intSize(v, strconv.IntSize)
	return int(x), err
}

// Int8 tries to parse Value as an int8 with strconv.ParseInt.
func (v Value) Int8() (int8, error) {
	x, err := intSize(v, 8)
	return int8(x), err
}

// Int16 tries to parse Value as an int16 with strconv.ParseInt.
func (v Value) Int16() (int16, error) {
	x, err := intSize(v, 16)
	return int16(x), err
}

// Int32 tries to parse Value as an int32 with strconv.ParseInt.
func (v Value) Int32() (int32, error) {
	x, err := intSize(v, 32)
	return int32(x), err
}

// Int64 tries to parse Value as an int64 with strconv.ParseInt.
func (v Value) Int64() (int64, error) {
	return intSize(v, 64)
}

// Uint tries to parse Value as an uint with strconv.ParseUint.
func (v Value) Uint() (uint, error) {
	x, err := uintSize(v, strconv.IntSize)
	return uint(x), err
}

// Uint8 tries to parse Value as an uint8 with strconv.ParseUint.
func (v Value) Uint8() (uint8, error) {
	x, err := uintSize(v, 8)
	return uint8(x), err
}

// Uint16 tries to parse Value as an uint16 with strconv.ParseUint.
func (v Value) Uint16() (uint16, error) {
	x, err := uintSize(v, 16)
	return uint16(x), err
}

// Uint32 tries to parse Value as an uint32 with strconv.ParseUint.
func (v Value) Uint32() (uint32, error) {
	x, err := uintSize(v, 32)
	return uint32(x), err
}

// Uint64 tries to parse Value as an uint64 with strconv.ParseUint.
func (v Value) Uint64() (uint64, error) {
	return uintSize(v, 64)
}

// Float32 tries to parse Value as a float32 with strconv.ParseFloat.
func (v Value) Float32() (float32, error) {
	x, err := floatSize(v, 32)
	return float32(x), err
}

// Float64 tries to parse Value as a float64 with strconv.ParseFloat.
func (v Value) Float64() (float64, error) {
	return floatSize(v, 64)
}

// Complex64 tries to parse Value as a complex64 with strconv.ParseComplex.
func (v Value) Complex64() (complex64, error) {
	x, err := complexSize(v, 64)
	return complex64(x), err
}

// Complex128 tries to parse Value as a complex128 with strconv.ParseComplex.
func (v Value) Complex128() (complex128, error) {
	return complexSize(v, 128)
}

// Duration tries to parse Value as time.Duration with time.ParseDuration.
func (v Value) Duration() (time.Duration, error) {
	x, err := time.ParseDuration(string(v))
	return x, errors.WithKind(err, ParseError)
}

// Url tries to parse Value as an url.Url with url.ParseRequestURI.
func (v Value) Url() (*url.URL, error) {
	x, err := url.ParseRequestURI(string(v))
	if err != nil {
		err = errors.WithKind(err, ParseError)
	}
	return x, err
}

func (v Value) UnmarshalText(u encoding.TextUnmarshaler) error {
	err := u.UnmarshalText([]byte(v))
	return errors.WithKind(err, errKind(err))
}

func intSize(v Value, bitSize int) (int64, error) {
	x, err := strconv.ParseInt(string(v), 0, bitSize)
	return x, errors.WithKind(err, errKind(err))
}

func uintSize(v Value, bitSize int) (uint64, error) {
	x, err := strconv.ParseUint(string(v), 0, bitSize)
	return x, errors.WithKind(err, errKind(err))
}

func floatSize(v Value, bitSize int) (float64, error) {
	x, err := strconv.ParseFloat(string(v), bitSize)
	return x, errors.WithKind(err, errKind(err))
}

func complexSize(v Value, bitSize int) (complex128, error) {
	x, err := strconv.ParseComplex(string(v), bitSize)
	return x, errors.WithKind(err, errKind(err))
}

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
