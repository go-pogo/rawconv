// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

const (
	DefaultItemsSeparator    = ","
	DefaultKeyValueSeparator = "="
)

type Options struct {
	ItemsSeparator    string // ,
	KeyValueSeparator string // =
}

func (o Options) itemSeparator() string {
	if o.ItemsSeparator == "" {
		return DefaultItemsSeparator
	}
	return o.ItemsSeparator
}

func (o Options) keyValueSeparator() string {
	if o.KeyValueSeparator == "" {
		return DefaultKeyValueSeparator
	}
	return o.KeyValueSeparator
}
