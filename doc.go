// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package rawconv implements conversions to and from raw string representations
of any (custom) data type in Go.

# Basic conversions

Out of the box, this package supports all logical base types and some common types:
- string
- bool
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- complex64, complex128
- array, slice
- map
- time.Duration
- url.URL
- encoding.TextUnmarshaler

# Array, slice and map conversions

Conversions to array, slice or map are done by splitting the raw string. The
separator can be set via the Options type and defaults to DefaultItemsSeparator.
For maps there is also a separator for the key-value pairs, which defaults to
DefaultKeyValueSeparator.
Values within the array, slice, or map are unmarshaled using the called
Unmarshaler. This is also done for keys of maps.
> Nested arrays, slices and maps are not supported.

# Structs

This package does not contain any logic for traversing struct types, because the
implementation would really depend on the use case. However, it is possible to
incorporate this package in your own struct unmarshaling logic.

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
