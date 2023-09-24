// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	nilChan := chan struct{}(nil)

	tests := map[string][]struct {
		input   any
		want    Value
		wantErr error
	}{
		"unsupported": {
			{
				input:   nilChan,
				wantErr: &UnsupportedTypeError{Type: reflect.TypeOf(nilChan)},
			},
		},

		"string": {
			{
				input: "test",
				want:  Value("test"),
			},
			{
				input: ptr("test"),
				want:  Value("test"),
			},
			{
				input: (*string)(nil),
			},
		},
		"bool": {
			{
				input: true,
				want:  Value("true"),
			},
			{
				input: ptr(false),
				want:  Value("false"),
			},
		},
		"int": {{
			input: -123,
			want:  Value("-123"),
		}},
		"uint": {{
			input: uint(123),
			want:  Value("123"),
		}},
		"float": {{
			input: 123.456,
			want:  Value("123.456"),
		}},
		"complex": {{
			input: 123.456 + 789.123i,
			want:  Value("(123.456+789.123i)"),
		}},

		"duration": {{
			input: (time.Second * 70) + (time.Millisecond * 123),
			want:  Value("1m10.123s"),
		}},
		"time": {{
			input: time.Date(1997, 8, 29, 13, 37, 1, 0, time.UTC),
			want:  Value("1997-08-29T13:37:01Z"),
		}},

		"url": {
			{
				input: &url.URL{Scheme: "http", Host: "localhost"},
				want:  Value("http://localhost"),
			},
			{
				input: (*url.URL)(nil),
				want:  Value(""),
			},
		},

		"ip": {{
			input: net.IPv4(192, 168, 1, 1),
			want:  Value("192.168.1.1"),
		}},
	}

	for name, tt := range tests {
		for _, tc := range tt {
			t.Run(name, func(t *testing.T) {
				have, haveErr := Marshal(tc.input)
				assert.Equal(t, tc.want, have)

				if tc.wantErr != nil {
					assert.ErrorIs(t, haveErr, tc.wantErr)
				} else {
					assert.NoError(t, haveErr)
				}
			})
		}
	}
}

func TestMarshalFunc_Exec(t *testing.T) {
	wantErr := errors.New("some err")
	_, haveErr := MarshalFunc(func(v any) (string, error) {
		return "", wantErr
	}).Exec(reflect.ValueOf("some value"))

	assert.ErrorIs(t, haveErr, wantErr)
}
