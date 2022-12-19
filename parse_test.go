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

func TestParser_Parse(t *testing.T) {
	ps := "some value"
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
		"*string": {
			val:  "some value",
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
		"*url": {
			val:  "http://localhost/",
			want: pu,
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
