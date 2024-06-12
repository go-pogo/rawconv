// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"
)

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
			assert.ErrorIs(t, err, ErrPointerExpected)
		})
	}
}

func TestUnmarshaler_Func(t *testing.T) {
	var u Unmarshaler
	u.Register(reflect.TypeOf(t), func(Value, any) error {
		return nil
	})
	testRegisterFind(t, 0, func(typ reflect.Type) any { return u.Func(typ) })
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
		"string": {{
			input: "some value",
			want:  "some value",
		}, {
			input: "foobar",
			want:  ptr("foobar"),
		}},
		"rune": {{
			input: "a",
			want:  'a',
		}, {
			input:   "abc",
			want:    'a',
			wantErr: ErrRuneTooManyChars,
		}},
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
		"time": {{
			input: "1997-08-29T13:37:00Z",
			want:  timeVal,
		}, {
			input: "1997-08-29T13:37:00Z",
			want:  ptr(timeVal),
		}},
		"url": {{
			input: "http://localhost/",
			want:  *urlPtr,
		}, {
			input: "http://localhost/",
			want:  urlPtr,
		}},
		"ip": {{
			input: "192.168.1.1",
			want:  net.IPv4(192, 168, 1, 1),
		}},
		"array": {{
			input: "1,2,3",
			want:  [3]int{1, 2, 3},
		}, {
			input: "1,2,3",
			want:  [3]string{"1", "2", "3"},
		}, {
			input:   "1,2,3",
			want:    [1]int{1},
			wantErr: ErrArrayTooManyValues,
		}, {
			input:   "1,2,3",
			want:    [1][2]int{},
			wantErr: ErrUnmarshalNested,
		}},
		"slice": {{
			input: "http://localhost/",
			want:  []url.URL{*urlPtr},
		}, {
			input: "1.2, 3.14, 5.6",
			want:  []float64{1.2, 3.14, 5.6},
		}, {
			input:   "iets",
			want:    ([][]string)(nil),
			wantErr: ErrUnmarshalNested,
		}},
		"map": {{
			input: "key1=value1,key2=value2",
			want:  map[string]string{"key1": "value1", "key2": "value2"},
		}, {
			input:   "iets",
			want:    map[string]map[string]string{},
			wantErr: ErrMapInvalidFormat,
		}, {
			input:   "iets=something",
			want:    map[string]map[string]string{},
			wantErr: ErrUnmarshalNested,
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
		assert.ErrorIs(t, err, ErrUnableToAddr)
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
