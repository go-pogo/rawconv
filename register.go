// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rawconv

import (
	"encoding"
	"net/url"
	"reflect"
	"time"
)

// RegisterUnmarshalFunc registers the UnmarshalFunc for typ, making it globally
// available for Unmarshal and any Unmarshaler.
func RegisterUnmarshalFunc(typ reflect.Type, fn UnmarshalFunc) {
	unmarshaler.Register(typ, fn)
}

// RegisterMarshalFunc registers the MarshalFunc for typ, making it globally
// available for Marshal, MarshalValue, MarshalReflect and any Marshaler.
func RegisterMarshalFunc(typ reflect.Type, fn MarshalFunc) {
	marshaler.Register(typ, fn)
}

func init() {
	// interfaces
	textMarshaler := reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	RegisterUnmarshalFunc(textMarshaler, unmarshalText)
	RegisterMarshalFunc(textMarshaler, marshalText)

	// common types
	rune := reflect.TypeOf(rune(0))
	RegisterUnmarshalFunc(rune, unmarshalRune)
	RegisterMarshalFunc(rune, marshalRune)

	timeDuration := reflect.TypeOf(time.Nanosecond)
	RegisterUnmarshalFunc(timeDuration, unmarshalDuration)
	RegisterMarshalFunc(timeDuration, marshalDuration)

	urlUrl := reflect.TypeOf(url.URL{})
	RegisterUnmarshalFunc(urlUrl, unmarshalUrl)
	RegisterMarshalFunc(urlUrl, marshalUrl)
}

func unmarshalText(val Value, dest any) error {
	return dest.(encoding.TextUnmarshaler).UnmarshalText(val.Bytes())
}

func marshalText(v any) (string, error) {
	b, err := v.(encoding.TextMarshaler).MarshalText()
	return string(b), err
}

type register[T interface{ MarshalFunc | UnmarshalFunc }] struct {
	types map[reflect.Kind]map[reflect.Type]int
	funcs []T
}

func (r *register[T]) initialized() bool { return r.types != nil && r.funcs != nil }

const panicUnsupportedKind = "rawconv: unsupported kind"

func (r *register[T]) add(typ reflect.Type, fn T) {
	k := typ.Kind()
	if k == reflect.Invalid ||
		k == reflect.Uintptr ||
		k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.UnsafePointer ||
		// not yet supported
		k == reflect.Array || k == reflect.Map || k == reflect.Slice {
		panic(panicUnsupportedKind)
	}

	// lazy init
	if r.types == nil {
		r.types = make(map[reflect.Kind]map[reflect.Type]int, 3)
	}
	if r.funcs == nil {
		r.funcs = make([]T, 0, 3)
	}

	if _, ok := r.types[k]; !ok {
		r.types[k] = map[reflect.Type]int{typ: len(r.funcs)}
	} else {
		r.types[k][typ] = len(r.funcs)
	}

	// store func
	r.funcs = append(r.funcs, fn)
}

func (r *register[T]) find(typ reflect.Type) T {
	// check if the exact type is registered
	if fn := r.getFromType(typ); fn != nil {
		return fn
	}

	if typ.Kind() != reflect.Ptr {
		// check if the type is registered as a pointer
		return r.getFromImpl(reflect.New(typ).Type())
	}

	// check if the elem type which is pointed to is registered
	if fn := r.find(typ.Elem()); fn != nil {
		return fn
	}
	if fn := r.getFromImpl(typ); fn != nil {
		return fn
	}

	return nil
}

func (r *register[T]) getFromType(typ reflect.Type) T {
	if kind, ok := r.types[typ.Kind()]; ok {
		if i, ok := kind[typ]; ok {
			return r.getFromIndex(i)
		}
	}
	return nil
}

func (r *register[T]) getFromImpl(typ reflect.Type) T {
	for x, i := range r.types[reflect.Interface] {
		if typ.Implements(x) {
			return r.getFromIndex(i)
		}
	}
	return nil
}

const panicInvalidFuncIndex = "rawconv: invalid index, func must exist!"

func (r *register[T]) getFromIndex(i int) T {
	if i >= len(r.funcs) {
		panic(panicInvalidFuncIndex)
	}
	return r.funcs[i]
}
