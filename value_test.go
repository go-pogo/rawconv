// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"math"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	types := map[string][]string{
		"empty":    {""},
		"string":   {"some value", "another string"},
		"rune":     {"a", "b", "c"},
		"bool":     {"true", "false"},
		"int":      {"100", "+33", "-349"},
		"float":    {"1.1", "0.59999", "22.564856"},
		"duration": {"10s", "2h13m12s"},
		"url":      {"https://foo.bar", "ftp://user@qux.xoo"},
	}

	type prepareFunc func(s string) (any, error)

	tests := map[string][3]prepareFunc{
		"String": {
			func(s string) (any, error) { return s, nil },
			func(s string) (any, error) { return Value(s).String(), nil },
			func(s string) (any, error) {
				var v string
				Value(s).StringVar(&v)
				return v, nil
			},
		},
		"Rune": {
			func(s string) (any, error) {
				var r rune
				if s != "" {
					r = rune(s[0])
				}
				return r, nil
			},
			func(s string) (any, error) { return Value(s).Rune(), nil },
			func(s string) (any, error) {
				var v rune
				Value(s).RuneVar(&v)
				return v, nil
			},
		},
		"Bytes": {
			func(s string) (any, error) { return []byte(s), nil },
			func(s string) (any, error) { return Value(s).Bytes(), nil },
			func(s string) (any, error) {
				var v []byte
				Value(s).BytesVar(&v)
				return v, nil
			},
		},
		"Bool": {
			func(s string) (any, error) { return strconv.ParseBool(s) },
			func(s string) (any, error) { return Value(s).Bool() },
			func(s string) (any, error) {
				var v bool
				err := Value(s).BoolVar(&v)
				return v, err
			},
		},
		"Int": {
			func(s string) (any, error) {
				i, err := strconv.ParseInt(s, 0, strconv.IntSize)
				return int(i), err
			},
			func(s string) (any, error) { return Value(s).Int() },
			func(s string) (any, error) {
				var v int
				err := Value(s).IntVar(&v)
				return v, err
			},
		},
		"Int8": {
			func(s string) (any, error) {
				i, err := strconv.ParseInt(s, 0, 8)
				return int8(i), err
			},
			func(s string) (any, error) { return Value(s).Int8() },
			func(s string) (any, error) {
				var v int8
				err := Value(s).Int8Var(&v)
				return v, err
			},
		},
		"Int16": {
			func(s string) (any, error) {
				i, err := strconv.ParseInt(s, 0, 16)
				return int16(i), err
			},
			func(s string) (any, error) { return Value(s).Int16() },
			func(s string) (any, error) {
				var v int16
				err := Value(s).Int16Var(&v)
				return v, err
			},
		},
		"Int32": {
			func(s string) (any, error) {
				i, err := strconv.ParseInt(s, 0, 32)
				return int32(i), err
			},
			func(s string) (any, error) { return Value(s).Int32() },
			func(s string) (any, error) {
				var v int32
				err := Value(s).Int32Var(&v)
				return v, err
			},
		},
		"Int64": {
			func(s string) (any, error) { return strconv.ParseInt(s, 0, 64) },
			func(s string) (any, error) { return Value(s).Int64() },
			func(s string) (any, error) {
				var v int64
				err := Value(s).Int64Var(&v)
				return v, err
			},
		},
		"Uint": {
			func(s string) (any, error) {
				i, err := strconv.ParseUint(s, 0, strconv.IntSize)
				return uint(i), err
			},
			func(s string) (any, error) { return Value(s).Uint() },
			func(s string) (any, error) {
				var v uint
				err := Value(s).UintVar(&v)
				return v, err
			},
		},
		"Uint8": {
			func(s string) (any, error) {
				i, err := strconv.ParseUint(s, 0, 8)
				return uint8(i), err
			},
			func(s string) (any, error) { return Value(s).Uint8() },
			func(s string) (any, error) {
				var v uint8
				err := Value(s).Uint8Var(&v)
				return v, err
			},
		},
		"Uint16": {
			func(s string) (any, error) {
				i, err := strconv.ParseUint(s, 0, 16)
				return uint16(i), err
			},
			func(s string) (any, error) { return Value(s).Uint16() },
			func(s string) (any, error) {
				var v uint16
				err := Value(s).Uint16Var(&v)
				return v, err
			},
		},
		"Uint32": {
			func(s string) (any, error) {
				i, err := strconv.ParseUint(s, 0, 32)
				return uint32(i), err
			},
			func(s string) (any, error) { return Value(s).Uint32() },
			func(s string) (any, error) {
				var v uint32
				err := Value(s).Uint32Var(&v)
				return v, err
			},
		},
		"Uint64": {
			func(s string) (any, error) { return strconv.ParseUint(s, 0, 64) },
			func(s string) (any, error) { return Value(s).Uint64() },
			func(s string) (any, error) {
				var v uint64
				err := Value(s).Uint64Var(&v)
				return v, err
			},
		},
		"Float32": {
			func(s string) (any, error) {
				i, err := strconv.ParseFloat(s, 32)
				return float32(i), err
			},
			func(s string) (any, error) { return Value(s).Float32() },
			func(s string) (any, error) {
				var v float32
				err := Value(s).Float32Var(&v)
				return v, err
			},
		},
		"Float64": {
			func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
			func(s string) (any, error) { return Value(s).Float64() },
			func(s string) (any, error) {
				var v float64
				err := Value(s).Float64Var(&v)
				return v, err
			},
		},
		"Complex64": {
			func(s string) (any, error) {
				i, err := strconv.ParseComplex(s, 64)
				return complex64(i), err
			},
			func(s string) (any, error) { return Value(s).Complex64() },
			func(s string) (any, error) {
				var v complex64
				err := Value(s).Complex64Var(&v)
				return v, err
			},
		},
		"Complex128": {
			func(s string) (any, error) { return strconv.ParseComplex(s, 128) },
			func(s string) (any, error) { return Value(s).Complex128() },
			func(s string) (any, error) {
				var v complex128
				err := Value(s).Complex128Var(&v)
				return v, err
			},
		},
		"Duration": {
			func(s string) (any, error) { return time.ParseDuration(s) },
			func(s string) (any, error) { return Value(s).Duration() },
			func(s string) (any, error) {
				var v time.Duration
				err := Value(s).DurationVar(&v)
				return v, err
			},
		},
		"Url": {
			func(s string) (any, error) { return url.ParseRequestURI(s) },
			func(s string) (any, error) { return Value(s).Url() },
			func(s string) (any, error) {
				var v url.URL
				err := Value(s).UrlVar(&v)

				if (v == url.URL{}) {
					return (*url.URL)(nil), err
				} else {
					return &v, err
				}
			},
		},
	}

	run := func(t *testing.T, prepWantFn, prepHaveFn prepareFunc) {
		for typ, inputs := range types {
			for _, input := range inputs {
				t.Run(typ, func(t *testing.T) {
					wantVal, wantErr := prepWantFn(input)
					haveVal, haveErr := prepHaveFn(input)

					assert.Exactlyf(t, wantVal, haveVal, "in: `%s`", input)
					assert.Exactly(t, wantErr, errors.Unwrap(haveErr))

					if wantErr != nil {
						assert.True(t,
							errors.Is(haveErr, ErrParseFailure) || errors.Is(haveErr, ErrValidationFailure),
							"err should be wrapped with either ErrParseFailure or ErrValidationFailure",
						)
					}
				})
			}
		}
	}

	for name, tcFn := range tests {
		t.Run(name, func(t *testing.T) { run(t, tcFn[0], tcFn[1]) })
		t.Run(name+"Var", func(t *testing.T) { run(t, tcFn[0], tcFn[2]) })
	}
}

func TestValue_IsEmpty(t *testing.T) {
	assert.True(t, Value("").IsEmpty())
	assert.False(t, Value("0").IsEmpty())
}

func TestValue_GoString(t *testing.T) {
	assert.Equal(t, `rawconv.Value("")`, Value("").GoString())
	assert.Equal(t, `rawconv.Value("0")`, Value("0").GoString())
	assert.Equal(t, `rawconv.Value("just some value")`, Value("just some value").GoString())
}

func TestValueFromComplex64(t *testing.T) {
	var want complex64 = 1 + 2i
	have, haveErr := ValueFromComplex64(want).Complex64()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromFloat32(t *testing.T) {
	var want float32 = math.Pi
	have, haveErr := ValueFromFloat32(want).Float32()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromInt(t *testing.T) {
	want := math.MinInt
	have, haveErr := ValueFromInt(want).Int()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromInt8(t *testing.T) {
	var want int8 = math.MaxInt8
	have, haveErr := ValueFromInt8(want).Int8()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromInt16(t *testing.T) {
	var want int16 = math.MinInt16
	have, haveErr := ValueFromInt16(want).Int16()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromInt32(t *testing.T) {
	var want int32 = math.MaxInt32
	have, haveErr := ValueFromInt32(want).Int32()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromUint(t *testing.T) {
	var want uint = math.MaxUint
	have, haveErr := ValueFromUint(want).Uint()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromUint8(t *testing.T) {
	var want uint8 = math.MaxUint8
	have, haveErr := ValueFromUint8(want).Uint8()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromUint16(t *testing.T) {
	var want uint16 = math.MaxUint16
	have, haveErr := ValueFromUint16(want).Uint16()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}

func TestValueFromUint32(t *testing.T) {
	var want uint32 = math.MaxUint32
	have, haveErr := ValueFromUint32(want).Uint32()
	assert.Equal(t, want, have)
	assert.Nil(t, haveErr)
}
