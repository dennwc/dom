//+build wasm

package js

import (
	"syscall/js"
)

var (
	global    = js.Global()
	null      = js.Null()
	undefined = js.Undefined()
)

// Ref is an alias for syscall/js.Value.
type Ref = js.Value

// Error is an alias for syscall/js.Error.
type Error = js.Error

// Type is a type name of a JS value, as returned by "typeof".
type Type = js.Type

const (
	TypeObject   = js.TypeObject
	TypeFunction = js.TypeFunction
)

// Wrapper is an alias for syscall/js.Wrapper.
type Wrapper = js.Wrapper

func valueOf(o interface{}) Ref {
	return js.ValueOf(toJS(o))
}

// Func is a wrapped Go function to be called by JavaScript.
type Func = js.Func

func funcOf(fnc func(this Ref, refs []Ref) interface{}) Func {
	return js.FuncOf(fnc)
}

func typedArrayOf(slice interface{}) Ref {
	return js.TypedArrayOf(slice).Value
}

func releaseTypedArray(v Ref) {
	js.TypedArray{v}.Release()
}
