// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"strconv"

	"github.com/go-pogo/errors"
)

// Int tries to parse Value as an int with strconv.ParseInt.
func (v Value) Int() (int, error) {
	x, err := intSize(v, strconv.IntSize)
	return int(x), err
}

// IntVar sets the value p points to using Int.
func (v Value) IntVar(p *int) (err error) {
	*p, err = v.Int()
	return
}

// Int8 tries to parse Value as an int8 with strconv.ParseInt.
func (v Value) Int8() (int8, error) {
	x, err := intSize(v, 8)
	return int8(x), err
}

// Int8Var sets the value p points to using Int8.
func (v Value) Int8Var(p *int8) (err error) {
	*p, err = v.Int8()
	return
}

// Int16 tries to parse Value as an int16 with strconv.ParseInt.
func (v Value) Int16() (int16, error) {
	x, err := intSize(v, 16)
	return int16(x), err
}

// Int16Var sets the value p points to using Int16.
func (v Value) Int16Var(p *int16) (err error) {
	*p, err = v.Int16()
	return
}

// Int32 tries to parse Value as an int32 with strconv.ParseInt.
func (v Value) Int32() (int32, error) {
	x, err := intSize(v, 32)
	return int32(x), err
}

// Int32Var sets the value p points to using Int32.
func (v Value) Int32Var(p *int32) (err error) {
	*p, err = v.Int32()
	return
}

// Int64 tries to parse Value as an int64 with strconv.ParseInt.
func (v Value) Int64() (int64, error) {
	return intSize(v, 64)
}

// Int64Var sets the value p points to using Int64.
func (v Value) Int64Var(p *int64) (err error) {
	*p, err = v.Int64()
	return
}

func intSize(v Value, bitSize int) (int64, error) {
	x, err := strconv.ParseInt(v.String(), 0, bitSize)
	return x, errors.WithKind(err, errKind(err))
}
