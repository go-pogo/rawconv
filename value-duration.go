// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"github.com/go-pogo/errors"
	"time"
)

// Duration tries to parse Value as a time.Duration using time.ParseDuration.
func (v Value) Duration() (time.Duration, error) {
	x, err := time.ParseDuration(v.String())
	return x, errors.Wrap(err, ErrParseFailure)
}

// DurationVar sets the value p points to using Duration.
func (v Value) DurationVar(p *time.Duration) (err error) {
	*p, err = v.Duration()
	return
}

func unmarshalDuration(val Value, dest any) error {
	if val.IsEmpty() {
		return nil
	}

	return val.DurationVar(dest.(*time.Duration))
}

func marshalDuration(v any) (string, error) {
	return v.(time.Duration).String(), nil
}
