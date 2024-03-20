rawconv
=======
[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/rawconv.svg?label=latest

[latest-release-url]: https://github.com/go-pogo/rawconv/releases

[build-status-img]: https://github.com/go-pogo/rawconv/actions/workflows/test.yml/badge.svg

[build-status-url]: https://github.com/go-pogo/rawconv/actions/workflows/test.yml

[report-img]: https://goreportcard.com/badge/github.com/go-pogo/rawconv

[report-url]: https://goreportcard.com/report/github.com/go-pogo/rawconv

[doc-img]: https://godoc.org/github.com/go-pogo/rawconv?status.svg

[doc-url]: https://pkg.go.dev/github.com/go-pogo/rawconv


Package `rawconv` implements conversions to and from raw string representations
of any (custom) data types in Go.

```sh
go get github.com/go-pogo/rawconv
```

```go
import "github.com/go-pogo/rawconv"
```

## Key features

- Convert from raw string to out of the box supported types:
  * `string`
  * `bool`
  * `int`, `int8`, `int16`, `int32`, `int64`
  * `uint`, `uint8`, `uint16`, `uint32`, `uint64`
  * `float32`, `float64`
  * `complex64`, `complex128`
  * `time.Duration`
  * `url.URL`
  * `encoding.TextUnmarshaler`
- Globally add support for your own custom types
- Or isolate support for your own custom types via `Marshaler` and `Unmarshaler` instances

```go
var duration time.Duration
if err := rawconv.Unmarshal("1h2m3s", &duration); err != nil {
    return fmt.Errorf("invalid duration: %w", err)
}
```

## Array, slice and map conversions

Conversions to `array`, `slice` or `map` are done by splitting the raw string. The separator can be set via the 
`Options` type and defaults to `DefaultItemsSeparator`. For maps there is also a separator for the key-value pairs, 
which defaults to `DefaultKeyValueSeparator`.
Values within the `array`, `slice`, or `map` are unmarshaled using the called `Unmarshaler`. This is also done for keys 
of maps.
> Nested arrays, slices and maps are not supported.

## Structs

This package does not contain any logic for traversing `struct` types, because the implementation would really depend 
on the use case. However, it is possible to incorporate this package in your own struct unmarshaling logic.

## Custom types

Custom types are supported in two ways; by implementing the `encoding.TextUnmarshaler` and/or `encoding.TextMarshaler`
interfaces, or by registering a `MarshalFunc` with `RegisterMarshalFunc` and/or an `UnmarshalFunc` with
`RegisterUnmarshalFunc`.
If you do not wish to globally expose your `MarshalFunc` or`UnmarshalFunc` implementations, it is possible to register
them to a new `Marshaler` or `Unmarshaler` and use those instances in your application instead.

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with

<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License

Copyright © 2022-2024 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
