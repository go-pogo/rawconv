// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package rawconv implements conversions to and from raw string representations
of any (custom) data type in Go.

# Basic conversions

Basic conversions are done using the strconv package. The following conversions
are supported by default:
- string
- bool
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- complex64, complex128
- time.Duration
- url.URL
These types are also implemented as methods on the Value type.

# Array, slice and map conversions

# Custom types

Custom types are supported in two ways; by implementing the
encoding.TextUnmarshaler and/or encoding.TextMarshaler interfaces, or by
registering a MarshalFunc with RegisterMarshalFunc and/or an UnmarshalFunc with
RegisterUnmarshalFunc.
If you do not wish to globally expose your MarshalFunc or UnmarshalFunc
implementations, it is possible to register them to a new Marshaler and/or
Unmarshaler and use those instances in your application instead.
*/
package rawconv
