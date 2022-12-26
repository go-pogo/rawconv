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

func TestUnmarshal(t *testing.T) {
	t.Run("ptr expected", func(t *testing.T) {
		var x string
		assert.ErrorIs(t, Unmarshal("some value", x), ErrPointerExpected)
	})
	t.Run("string", func(t *testing.T) {
		var have string
		want := "test"
		assert.Nil(t, Unmarshal(Value(want), &have))
		assert.Equal(t, want, have)
	})
	t.Run("url", func(t *testing.T) {
		var have url.URL
		want, err := url.ParseRequestURI("https://example.com:1234/somepath?xyz")
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, Unmarshal(Value(want.String()), &have))
		assert.Equal(t, want, &have)
	})

}

func TestParser_Parse(t *testing.T) {
	pt, _ := time.Parse(time.RFC3339, "1997-08-29T13:37:00Z")
	pu, _ := url.ParseRequestURI("http://localhost/")

	tests := map[string]struct {
		val     Value
		want    interface{}
		wantErr error
	}{
		"string": {
			val:  "some value",
			want: "some value",
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
			rtyp := reflect.TypeOf(tc.want)
			// reflect.New returns a pointer to the type,
			// Elem() removes that pointer
			rval := reflect.New(rtyp).Elem()

			haveErr := Parse(tc.val, rval)
			assert.Equal(t, tc.want, rval.Interface())

			if tc.wantErr != nil {
				assert.ErrorIs(t, haveErr, tc.wantErr)
			} else {
				assert.Nil(t, haveErr)
			}
		})
	}
}
