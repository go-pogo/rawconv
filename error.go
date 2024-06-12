// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"github.com/go-pogo/errors"
	"reflect"
	"strconv"
)

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Is(err error) bool {
	//goland:noinspection GoTypeAssertionOnErrors
	t, ok := err.(*UnsupportedTypeError)
	return ok && e.Type == t.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "type `" + e.Type.String() + "` is not supported"
}

const (
	ErrParseFailure      errors.Msg = "failed to parse"
	ErrValidationFailure errors.Msg = "failed to validate"
)

func errKind(err error) error {
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		if errors.Is(numErr.Err, strconv.ErrRange) {
			return ErrValidationFailure
		} else {
			return ErrParseFailure
		}
	}
	return nil
}
