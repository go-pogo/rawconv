// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"github.com/go-pogo/errors"
	"time"
)

func (v Value) Duration() (time.Duration, error) {
	x, err := time.ParseDuration(v.String())
	return x, errors.WithKind(err, ParseError)
}

func (v Value) DurationVar(p *time.Duration) (err error) {
	*p, err = v.Duration()
	return
}

func parseDuration(val Value, dest interface{}) error {
	if val.Empty() {
		return nil
	}

	return val.DurationVar(dest.(*time.Duration))
}
