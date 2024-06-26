// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"strconv"

	"github.com/go-pogo/errors"
)

// ValueFromBool encodes v to a Value using strconv.FormatBool.
func ValueFromBool(v bool) Value {
	return Value(strconv.FormatBool(v))
}

// Bool tries to parse Value as a bool using strconv.ParseBool.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns an error.
func (v Value) Bool() (bool, error) {
	x, err := strconv.ParseBool(string(v))
	if kind := errKind(err); kind != nil {
		return x, errors.Wrap(err, kind)
	}
	return x, errors.WithStack(err)
}

// BoolVar sets the value p points to using Bool.
func (v Value) BoolVar(p *bool) (err error) {
	*p, err = v.Bool()
	return
}
