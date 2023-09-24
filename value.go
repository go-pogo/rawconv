// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

// Value is a textual representation of a raw value which is able to cast itself
// to any of the supported types using its corresponding method.
//
//	boolVal, err := rawconv.Value("true").Bool()
type Value string

// IsEmpty indicates if Value is an empty string.
func (v Value) IsEmpty() bool { return string(v) == "" }

func (v Value) GoString() string { return `rawconv.Value("` + v.String() + `")` }

// String returns Value as a raw string.
func (v Value) String() string { return string(v) }

// StringVar sets the value p points to, to Value as raw string.
func (v Value) StringVar(p *string) { *p = v.String() }

// Bytes returns Value as raw bytes.
func (v Value) Bytes() []byte { return []byte(v) }

// BytesVar sets the value p points to, to Value as raw bytes.
func (v Value) BytesVar(p *[]byte) { *p = v.Bytes() }
