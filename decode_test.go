// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func assertInvalidActionError(t *testing.T, err, target error) {
	assert.ErrorIs(t, err, target)
	assert.ErrorIs(t, err, ImplementationError)
}

func ptr[T any](v T) *T { return &v }

func TestUnmarshal(t *testing.T) {
	tests := map[string]any{
		"nil":           nil,
		"not a pointer": struct{}{},
		"nil pointer":   (*url.URL)(nil),
	}
	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			err := Unmarshal("", val)
			assertInvalidActionError(t, err, ErrPointerExpected)
		})
	}
}

func TestUnmarshaler_Unmarshal(t *testing.T) {
	timeVal, _ := time.Parse(time.RFC3339, "1997-08-29T13:37:00Z")
	urlPtr, _ := url.ParseRequestURI("http://localhost/")
	nilChan := chan struct{}(nil)

	tests := map[string][]struct {
		input   Value
		want    any
		wantErr error
	}{
		"empty": {{
			input: "",
			want:  "",
		}},
		"unsupported": {{
			input:   "some value",
			want:    nilChan,
			wantErr: &UnsupportedTypeError{Type: reflect.TypeOf(&nilChan)},
		}},

		"string": {
			{
				input: "some value",
				want:  "some value",
			},
			{
				input: Value("foobar"),
				want:  ptr("foobar"),
			},
		},
		"bool": {{
			input: "true",
			want:  true,
		}},
		"int": {{
			input: "-10",
			want:  -10,
		}},
		"uint": {{
			input: "1337",
			want:  uint(1337),
		}},
		"float": {{
			input: "3.14",
			want:  3.14,
		}},
		"complex": {{
			input: "3.14+2.72i",
			want:  complex(3.14, 2.72),
		}},

		"duration": {{
			input: "10s",
			want:  time.Second * 10,
		}},
		"time": {
			{
				input: "1997-08-29T13:37:00Z",
				want:  timeVal,
			},
			{
				input: "1997-08-29T13:37:00Z",
				want:  ptr(timeVal),
			},
		},

		"url": {
			{
				input: "http://localhost/",
				want:  *urlPtr,
			},
			{
				input: "http://localhost/",
				want:  urlPtr,
			},
		},

		"ip": {{
			input: "192.168.1.1",
			want:  net.IPv4(192, 168, 1, 1),
		}},
	}

	for name, tt := range tests {
		for _, tc := range tt {
			t.Run(name, func(t *testing.T) {
				rt := reflect.TypeOf(tc.want)
				rv := reflect.New(rt)

				haveErr := unmarshaler.Unmarshal(tc.input, rv)
				assert.Equal(t, tc.want, rv.Elem().Interface())

				if tc.wantErr != nil {
					assert.ErrorIs(t, haveErr, tc.wantErr)
				} else {
					assert.NoError(t, haveErr)
				}
			})
		}
	}

	t.Run("unable to set", func(t *testing.T) {
		haveErr := unmarshaler.Unmarshal("some value", reflect.ValueOf("some value"))
		assert.ErrorIs(t, haveErr, ErrUnableToSet)
	})
}

func TestParseFunc_Exec(t *testing.T) {
	durationType := reflect.TypeOf(time.Nanosecond)
	parseFunc := UnmarshalFunc(unmarshalDuration)

	t.Parallel()
	t.Run("value value", func(t *testing.T) {
		err := parseFunc.Exec("10s", reflect.Zero(durationType))
		assertInvalidActionError(t, err, ErrUnableToAddr)
	})
	t.Run("nil pointer", func(t *testing.T) {
		var d *time.Duration
		assert.NoError(t, parseFunc.Exec("10s", reflect.ValueOf(&d)))
		assert.Equal(t, time.Second*10, *d)
	})
	t.Run("multiple pointers", func(t *testing.T) {
		var d ***time.Duration
		assert.NoError(t, parseFunc.Exec("10s", reflect.ValueOf(&d)))
		assert.Equal(t, time.Second*10, ***d)
	})
	t.Run("new", func(t *testing.T) {
		rv := reflect.New(durationType)
		assert.NoError(t, parseFunc.Exec("10s", rv))
		assert.Equal(t, time.Second*10, rv.Elem().Interface())
	})
	t.Run("slice index", func(t *testing.T) {
		slice := []time.Duration{time.Second}
		assert.NoError(t, parseFunc.Exec("10s", reflect.ValueOf(slice).Index(0)))
		assert.Equal(t, time.Second*10, slice[0])
	})
}
