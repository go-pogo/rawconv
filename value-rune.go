// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import "github.com/go-pogo/errors"

// Rune returns the first rune of Value.
func (v Value) Rune() rune {
	if v.IsEmpty() {
		return rune(0)
	}
	return rune(v[0])
}

// RuneVar sets the value p points to, to the first rune of Value.
func (v Value) RuneVar(p *rune) { *p = v.Rune() }

func unmarshalRune(val Value, dest any) error {
	val.RuneVar(dest.(*rune))
	if len(val) > 1 {
		return errors.New(ErrRuneTooManyChars)
	}
	return nil
}

func marshalRune(v any) (string, error) {
	return string(v.(rune)), nil
}
