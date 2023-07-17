// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/url"
	"reflect"
	"time"
)

func ExampleUnmarshal() {
	var duration time.Duration
	if err := Unmarshal("1h2m3s", &duration); err != nil {
		panic(err)
	}

	fmt.Println(duration)
	// Output: 1h2m3s
}

func ExampleParse() {
	var website *url.URL
	if err := UnmarshalReflect("https://example.com", reflect.ValueOf(&website)); err != nil {
		panic(err)
	}

	fmt.Println(website.String())
	// Output: https://example.com
}

func ExampleParser_Parse() {
	type myType struct {
		something string
	}

	var u Unmarshaler
	u.Register(reflect.TypeOf(myType{}), func(val Value, dest interface{}) error {
		mt := dest.(*myType)
		mt.something = val.String()
		return nil
	})

	var mt myType
	if err := u.Unmarshal("some value", reflect.ValueOf(&mt)); err != nil {
		panic(err)
	}

	spew.Dump(mt)
	// Output:
	// (parseval.myType) {
	//  something: (string) (len=10) "some value"
	// }
}
