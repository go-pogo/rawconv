// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parseval

import (
	"fmt"
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
	if err := Parse("http://example.com", reflect.ValueOf(&website)); err != nil {
		panic(err)
	}

	fmt.Println(website.String())
	// Output: http://example.com
}
