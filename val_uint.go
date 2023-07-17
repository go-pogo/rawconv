// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"strconv"

	"github.com/go-pogo/errors"
)

// Uint tries to parse Value as an uint using strconv.ParseUint
func (v Value) Uint() (uint, error) {
	x, err := uintSize(v, strconv.IntSize)
	return uint(x), err
}

// UintVar sets the value p points to using Uint.
func (v Value) UintVar(p *uint) (err error) {
	*p, err = v.Uint()
	return
}

// Uint8 tries to parse Value as an uint8 using strconv.ParseUint.
func (v Value) Uint8() (uint8, error) {
	x, err := uintSize(v, 8)
	return uint8(x), err
}

// Uint8Var sets the value p points to using Uint8.
func (v Value) Uint8Var(p *uint8) (err error) {
	*p, err = v.Uint8()
	return
}

// Uint16 tries to parse Value as an uint16 using strconv.ParseUint.
func (v Value) Uint16() (uint16, error) {
	x, err := uintSize(v, 16)
	return uint16(x), err
}

// Uint16Var sets the value p points to using Uint16.
func (v Value) Uint16Var(p *uint16) (err error) {
	*p, err = v.Uint16()
	return
}

// Uint32 tries to parse Value as an uint32 using strconv.ParseUint.
func (v Value) Uint32() (uint32, error) {
	x, err := uintSize(v, 32)
	return uint32(x), err
}

// Uint32Var sets the value p points to using Uint32.
func (v Value) Uint32Var(p *uint32) (err error) {
	*p, err = v.Uint32()
	return
}

// Uint64 tries to parse Value as an uint64 using strconv.ParseUint.
func (v Value) Uint64() (uint64, error) {
	return uintSize(v, 64)
}

// Uint64Var sets the value p points to using Uint64.
func (v Value) Uint64Var(p *uint64) (err error) {
	*p, err = v.Uint64()
	return
}

func uintSize(v Value, bitSize int) (uint64, error) {
	x, err := strconv.ParseUint(v.String(), 0, bitSize)
	return x, errors.WithKind(err, errKind(err))
}
