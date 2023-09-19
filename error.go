// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

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
	ParseError      errors.Kind = "parse error"
	ValidationError errors.Kind = "validation error"
)

func errKind(err error) errors.Kind {
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		if errors.Is(numErr.Err, strconv.ErrRange) {
			return ValidationError
		} else {
			return ParseError
		}
	}
	return errors.UnknownKind
}
