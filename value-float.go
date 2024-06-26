// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"strconv"

	"github.com/go-pogo/errors"
)

// ValueFromFloat32 encodes v to a Value using strconv.FormatFloat.
func ValueFromFloat32(v float32) Value {
	return Value(strconv.FormatFloat(float64(v), 'g', -1, 32))
}

// ValueFromFloat64 encodes v to a Value using strconv.FormatFloat.
func ValueFromFloat64(v float64) Value {
	return Value(strconv.FormatFloat(v, 'g', -1, 64))
}

// Float32 tries to parse Value as a float32 using strconv.ParseFloat.
func (v Value) Float32() (float32, error) {
	x, err := floatSize(v, 32)
	return float32(x), err
}

// Float32Var sets the value p points to using Float32.
func (v Value) Float32Var(p *float32) (err error) {
	*p, err = v.Float32()
	return
}

// Float64 tries to parse Value as a float64 using strconv.ParseFloat.
func (v Value) Float64() (float64, error) {
	return floatSize(v, 64)
}

// Float64Var sets the value p points to using Float64.
func (v Value) Float64Var(p *float64) (err error) {
	*p, err = v.Float64()
	return
}

func floatSize(v Value, bitSize int) (float64, error) {
	x, err := strconv.ParseFloat(v.String(), bitSize)
	if kind := errKind(err); kind != nil {
		return x, errors.Wrap(err, kind)
	}
	return x, errors.WithStack(err)
}
