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
	assert.ErrorIs(t, err, InvalidActionError)
}

func TestFunc(t *testing.T) {
	tests := []struct {
		want  uintptr
		types []reflect.Type
	}{
		{
			want: reflect.ValueOf(unmarshalText).Pointer(),
			types: []reflect.Type{
				reflect.TypeOf(net.IP{}),
				reflect.TypeOf((*net.IP)(nil)),
				reflect.TypeOf((**net.IP)(nil)),
			},
		},
		{
			want: reflect.ValueOf(unmarshalDuration).Pointer(),
			types: []reflect.Type{
				reflect.TypeOf(time.Second),
				reflect.TypeOf((*time.Duration)(nil)),
				reflect.TypeOf((**time.Duration)(nil)),
			},
		},
		{
			want: reflect.ValueOf(unmarshalUrl).Pointer(),
			types: []reflect.Type{
				reflect.TypeOf(url.URL{}),
				reflect.TypeOf((*url.URL)(nil)),
				reflect.TypeOf((**url.URL)(nil)),
			},
		},
	}

	for _, tc := range tests {
		for _, typ := range tc.types {
			t.Run(typ.String(), func(t *testing.T) {
				have := Func(typ)
				assert.Equal(t, tc.want, reflect.ValueOf(have).Pointer())
			})
		}
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	tests := map[string]interface{}{
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

	t.Run("string", func(t *testing.T) {
		var have string
		want := "test"
		assert.Nil(t, Unmarshal(Value(want), &have))
		assert.Equal(t, want, have)
	})
	t.Run("url", func(t *testing.T) {
		var have *url.URL
		want, err := url.ParseRequestURI("https://example.com:1234/somepath?xyz")
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, Unmarshal(Value(want.String()), &have))
		assert.Equal(t, want, have)
	})
	t.Run("empty url", func(t *testing.T) {
		var have *url.URL
		assert.Nil(t, Unmarshal("", &have))
		assert.Equal(t, new(url.URL), have)
	})
	t.Run("ip", func(t *testing.T) {
		var have *net.IP
		want := net.ParseIP("10.0.0.10")

		assert.Nil(t, Unmarshal(Value(want.String()), &have))
		assert.Equal(t, want, *have)
	})
}

func TestParse(t *testing.T) {
	pt, _ := time.Parse(time.RFC3339, "1997-08-29T13:37:00Z")
	pu, _ := url.ParseRequestURI("http://localhost/")
	ps := "pointer to string"

	tests := map[string]struct {
		val     Value
		want    interface{}
		wantErr error
	}{
		"string": {
			val:  "some value",
			want: "some value",
		},
		"*string": {
			val:  Value(ps),
			want: &ps,
		},

		"duration": {
			val:  "10s",
			want: time.Second * 10,
		},
		"time": {
			val:  "1997-08-29T13:37:00Z",
			want: pt,
		},

		"url": {
			val:  "http://localhost/",
			want: *pu,
		},

		"ip": {
			val:  "192.168.1.1",
			want: net.IPv4(192, 168, 1, 1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rt := reflect.TypeOf(tc.want)
			rv := reflect.New(rt)

			haveErr := UnmarshalReflect(tc.val, rv)
			assert.Equal(t, tc.want, rv.Elem().Interface())

			if tc.wantErr != nil {
				assert.ErrorIs(t, haveErr, tc.wantErr)
			} else {
				assert.Nil(t, haveErr)
			}
		})
	}
}

func TestParseFunc_Exec(t *testing.T) {
	typ := reflect.TypeOf(time.Nanosecond)
	parseFunc := UnmarshalFunc(unmarshalDuration)

	t.Parallel()
	t.Run("value value", func(t *testing.T) {
		err := parseFunc.Exec("10s", reflect.Zero(typ))
		assertInvalidActionError(t, err, ErrUnableToAddr)
	})
	t.Run("nil pointer", func(t *testing.T) {
		var d *time.Duration
		assert.Nil(t, parseFunc.Exec("10s", reflect.ValueOf(&d)))
		assert.Equal(t, time.Second*10, *d)
	})
	t.Run("multiple pointers", func(t *testing.T) {
		var d ***time.Duration
		assert.Nil(t, parseFunc.Exec("10s", reflect.ValueOf(&d)))
		assert.Equal(t, time.Second*10, ***d)
	})
	t.Run("new", func(t *testing.T) {
		rv := reflect.New(typ)
		assert.Nil(t, parseFunc.Exec("10s", rv))
		assert.Equal(t, time.Second*10, rv.Elem().Interface())
	})
}
