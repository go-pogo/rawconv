// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

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

func ExampleMarshal() {
	duration := time.Hour + (time.Minute * 2) + (time.Second * 3)
	val, err := Marshal(duration)
	if err != nil {
		panic(err)
	}

	fmt.Println(val.String())
	// Output: 1h2m3s
}

func ExampleUnmarshaler() {
	var u Unmarshaler
	var target *url.URL
	if err := u.Unmarshal("https://example.com", reflect.ValueOf(&target)); err != nil {
		panic(err)
	}

	fmt.Println(target.String())
	// Output: https://example.com
}

func ExampleMarshaler() {
	var m Marshaler
	target, _ := url.ParseRequestURI("https://example.com")
	val, err := m.Marshal(reflect.ValueOf(target))
	if err != nil {
		panic(err)
	}

	fmt.Println(val.String())
	// Output: https://example.com
}

func ExampleUnmarshaler_Register() {
	type myType struct {
		something string
	}

	var u Unmarshaler
	u.Register(reflect.TypeOf(myType{}), func(val Value, dest any) error {
		mt := dest.(*myType)
		mt.something = val.String()
		return nil
	})

	var target myType
	if err := u.Unmarshal("some value", reflect.ValueOf(&target)); err != nil {
		panic(err)
	}

	spew.Dump(target)
	// Output:
	// (rawconv.myType) {
	//  something: (string) (len=10) "some value"
	// }
}
