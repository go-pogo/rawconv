// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"github.com/go-pogo/errors"
	"net/url"
)

// Url tries to parse Value as an *url.URL using url.ParseRequestURI.
func (v Value) Url() (*url.URL, error) {
	x, err := url.ParseRequestURI(v.String())
	if err != nil {
		return nil, errors.WithKind(err, ParseError)
	}
	return x, nil
}

// UrlVar sets the value p points to using Url.
func (v Value) UrlVar(p *url.URL) error {
	x, err := v.Url()
	if err != nil {
		return err
	}
	*p = *x
	return nil
}

func unmarshalUrl(val Value, dest interface{}) error {
	if val.Empty() {
		return nil
	}

	return val.UrlVar(dest.(*url.URL))
}
