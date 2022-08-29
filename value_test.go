// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"net"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	types := map[string][]string{
		"empty":  {""},
		"string": {"some value", "another string"},
		"bool":   {"true", "false"},
		"int":    {"100", "+33", "-349"},
		"float":  {"1.1", "0.59999", "22.564856"},
		"time":   {"10s", "2h"},
		"url":    {"https://foo.bar", "user:pass@tcp(12.34.56.78:1337)/qux"},
	}

	tests := map[string][2]func(s string) (interface{}, error){
		"String": {
			func(s string) (interface{}, error) { return Value(s).String(), nil },
			func(s string) (interface{}, error) { return s, nil },
		},
		"Bool": {
			func(s string) (interface{}, error) { return Value(s).Bool() },
			func(s string) (interface{}, error) { return strconv.ParseBool(s) },
		},
		"Int": {
			func(s string) (interface{}, error) { return Value(s).Int() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, strconv.IntSize)
				return int(i), err
			},
		},
		"Int8": {
			func(s string) (interface{}, error) { return Value(s).Int8() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, 8)
				return int8(i), err
			},
		},
		"Int16": {
			func(s string) (interface{}, error) { return Value(s).Int16() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, 16)
				return int16(i), err
			},
		},
		"Int32": {
			func(s string) (interface{}, error) { return Value(s).Int32() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, 32)
				return int32(i), err
			},
		},
		"Int64": {
			func(s string) (interface{}, error) { return Value(s).Int64() },
			func(s string) (interface{}, error) { return strconv.ParseInt(s, 0, 64) },
		},
		"Uint": {
			func(s string) (interface{}, error) { return Value(s).Uint() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, strconv.IntSize)
				return uint(i), err
			},
		},
		"Uint8": {
			func(s string) (interface{}, error) { return Value(s).Uint8() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, 8)
				return uint8(i), err
			},
		},
		"Uint16": {
			func(s string) (interface{}, error) { return Value(s).Uint16() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, 16)
				return uint16(i), err
			},
		},
		"Uint32": {
			func(s string) (interface{}, error) { return Value(s).Uint32() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, 32)
				return uint32(i), err
			},
		},
		"Uint64": {
			func(s string) (interface{}, error) { return Value(s).Uint64() },
			func(s string) (interface{}, error) { return strconv.ParseUint(s, 0, 64) },
		},
		"Float32": {
			func(s string) (interface{}, error) { return Value(s).Float32() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseFloat(s, 32)
				return float32(i), err
			},
		},
		"Float64": {
			func(s string) (interface{}, error) { return Value(s).Float64() },
			func(s string) (interface{}, error) { return strconv.ParseFloat(s, 64) },
		},
		"Complex64": {
			func(s string) (interface{}, error) { return Value(s).Complex64() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseComplex(s, 64)
				return complex64(i), err
			},
		},
		"Complex128": {
			func(s string) (interface{}, error) { return Value(s).Complex128() },
			func(s string) (interface{}, error) { return strconv.ParseComplex(s, 128) },
		},
		"Duration": {
			func(s string) (interface{}, error) { return Value(s).Duration() },
			func(s string) (interface{}, error) { return time.ParseDuration(s) },
		},
		"Url": {
			func(s string) (interface{}, error) { return Value(s).Url() },
			func(s string) (interface{}, error) { return url.ParseRequestURI(s) },
		},
	}

	for name, tcFn := range tests {
		t.Run(name, func(t *testing.T) {
			for typ, inputs := range types {
				for _, input := range inputs {
					t.Run(typ, func(t *testing.T) {
						haveVal, haveErr := tcFn[0](input)
						wantVal, wantErr := tcFn[1](input)

						assert.Exactly(t, wantVal, haveVal)
						assert.Exactly(t, wantErr, errors.Unwrap(haveErr))

						if wantErr != nil {
							haveKind := errors.GetKind(haveErr)
							assert.True(t, haveKind == ParseError || haveKind == ValidationError, "Kind should match ParseError or ValidationError")
						}
					})
				}
			}
		})
	}
}

func TestValue_Empty(t *testing.T) {
	assert.True(t, Value("").Empty())
	assert.False(t, Value("0").Empty())
}

func TestValue_UnmarshalText(t *testing.T) {
	input := "127.0.0.38"

	var target net.IP
	assert.Nil(t, Value(input).UnmarshalText(&target))
	assert.Equal(t, net.ParseIP(input), target)
}
