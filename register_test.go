package parseval

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestGetUnmarshalFunc(t *testing.T) {
	testRegisterFind(t, 0, func(typ reflect.Type) any { return GetUnmarshalFunc(typ) })
}

func TestGetMarshalFunc(t *testing.T) {
	testRegisterFind(t, 1, func(typ reflect.Type) any { return GetMarshalFunc(typ) })
}

func testRegisterFind(t *testing.T, i int, haveFn func(reflect.Type) any) {
	tests := []struct {
		want  [2]uintptr
		types []reflect.Type
	}{
		{
			want: [2]uintptr{
				reflect.ValueOf(unmarshalText).Pointer(),
				reflect.ValueOf(marshalText).Pointer(),
			},
			types: []reflect.Type{
				reflect.TypeOf(net.IP{}),
				reflect.TypeOf((*net.IP)(nil)),
				reflect.TypeOf((**net.IP)(nil)),
			},
		},
		{
			want: [2]uintptr{
				reflect.ValueOf(unmarshalDuration).Pointer(),
				reflect.ValueOf(marshalDuration).Pointer(),
			},
			types: []reflect.Type{
				reflect.TypeOf(time.Second),
				reflect.TypeOf((*time.Duration)(nil)),
				reflect.TypeOf((**time.Duration)(nil)),
			},
		},
		{
			want: [2]uintptr{
				reflect.ValueOf(unmarshalUrl).Pointer(),
				reflect.ValueOf(marshalUrl).Pointer(),
			},
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
				have := haveFn(typ)
				assert.Equal(t, tc.want[i], reflect.ValueOf(have).Pointer())
			})
		}
	}
}
