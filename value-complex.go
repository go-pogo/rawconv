// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"strconv"

	"github.com/go-pogo/errors"
)

// ValueFromComplex64 encodes v to a Value using strconv.FormatComplex.
func ValueFromComplex64(v complex64) Value {
	return Value(strconv.FormatComplex(complex128(v), 'g', -1, 64))
}

// ValueFromComplex128 encodes v to a Value using strconv.FormatComplex.
func ValueFromComplex128(v complex128) Value {
	return Value(strconv.FormatComplex(v, 'g', -1, 128))
}

// Complex64 tries to parse Value as a complex64 using strconv.ParseComplex.
func (v Value) Complex64() (complex64, error) {
	x, err := complexSize(v, 64)
	return complex64(x), err
}

// Complex64Var sets the value p points to using Complex64.
func (v Value) Complex64Var(p *complex64) (err error) {
	*p, err = v.Complex64()
	return
}

// Complex128 tries to parse Value as a complex128 using strconv.ParseComplex.
func (v Value) Complex128() (complex128, error) {
	return complexSize(v, 128)
}

// Complex128Var sets the value p points to using Complex128.
func (v Value) Complex128Var(p *complex128) (err error) {
	*p, err = v.Complex128()
	return
}

func complexSize(v Value, bitSize int) (complex128, error) {
	x, err := strconv.ParseComplex(v.String(), bitSize)
	if kind := errKind(err); kind != nil {
		return x, errors.Wrap(err, kind)
	}
	return x, errors.WithStack(err)
}
