// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package rawconv implements conversions to and from raw string representations
of any (custom) data type in Go.

# Basic conversions

Basic conversions are done using the strconv package and are implemented as
methods on the Value type. The following conversions are supported:
- string
- bool
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- complex64, complex128
- time.Duration
- url.URL

# Custom types

Conversions for global custom types are done by registering a MarshalFunc and/or
UnmarshalFunc using the RegisterMarshalFunc and RegisterUnmarshalFunc functions.
It is also possible to use Marshaler and/or Unmarshaler if you do not want to
expose the MarshalFunc or UnmarshalFunc implementations.
*/
package rawconv
