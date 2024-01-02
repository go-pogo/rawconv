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

## Basic conversions

Basic conversions are done using the `strconv` package and are implemented as
methods on the `Value` type. The following conversions are supported:
- `string`
- `bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `complex64`, `complex128`
- `time.Duration`
- `url.URL`

## Custom types

Conversions for global custom types are done by registering a `MarshalFunc` and/or
`UnmarshalFunc` using the `RegisterMarshalFunc` and `RegisterUnmarshalFunc` functions.
It is also possible to use `Marshaler` and/or `Unmarshaler` if you do not want to
expose the `MarshalFunc` or `UnmarshalFunc` implementations.

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with

<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License

Copyright © 2022-2024 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
