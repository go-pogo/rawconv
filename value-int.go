// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"github.com/go-pogo/errors"
	"strconv"
)

// ValueFromInt encodes v to a Value using strconv.FormatInt.
func ValueFromInt(v int) Value {
	return Value(strconv.FormatInt(int64(v), 10))
}

// ValueFromInt8 encodes v to a Value using strconv.FormatInt.
func ValueFromInt8(v int8) Value {
	return Value(strconv.FormatInt(int64(v), 10))
}

// ValueFromInt16 encodes v to a Value using strconv.FormatInt.
func ValueFromInt16(v int16) Value {
	return Value(strconv.FormatInt(int64(v), 10))
}

// ValueFromInt32 encodes v to a Value using strconv.FormatInt.
func ValueFromInt32(v int32) Value {
	return Value(strconv.FormatInt(int64(v), 10))
}

// ValueFromInt64 encodes v to a Value using strconv.FormatInt.
func ValueFromInt64(v int64) Value {
	return Value(strconv.FormatInt(v, 10))
}

// Int tries to parse Value as an int using strconv.ParseInt.
func (v Value) Int() (int, error) {
	x, err := intSize(v, strconv.IntSize)
	return int(x), err
}

// IntVar sets the value p points to using Int.
func (v Value) IntVar(p *int) (err error) {
	*p, err = v.Int()
	return
}

// Int8 tries to parse Value as an int8 using strconv.ParseInt.
func (v Value) Int8() (int8, error) {
	x, err := intSize(v, 8)
	return int8(x), err
}

// Int8Var sets the value p points to using Int8.
func (v Value) Int8Var(p *int8) (err error) {
	*p, err = v.Int8()
	return
}

// Int16 tries to parse Value as an int16 using strconv.ParseInt.
func (v Value) Int16() (int16, error) {
	x, err := intSize(v, 16)
	return int16(x), err
}

// Int16Var sets the value p points to using Int16.
func (v Value) Int16Var(p *int16) (err error) {
	*p, err = v.Int16()
	return
}

// Int32 tries to parse Value as an int32 using strconv.ParseInt.
func (v Value) Int32() (int32, error) {
	x, err := intSize(v, 32)
	return int32(x), err
}

// Int32Var sets the value p points to using Int32.
func (v Value) Int32Var(p *int32) (err error) {
	*p, err = v.Int32()
	return
}

// Int64 tries to parse Value as an int64 using strconv.ParseInt.
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
	if kind := errKind(err); kind != nil {
		return x, errors.Wrap(err, kind)
	}
	return x, errors.WithStack(err)
}
