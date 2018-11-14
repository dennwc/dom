//+build wasm,!go1.12

package js

import "syscall/js"

func unwrap(v Wrapper) Ref {
	return v.JSValue()
}

// JSRef is a common interface for object that are backed by a JS object.
//
// Deprecated: see Wrapper
type JSRef = Wrapper

// Wrapper is implemented by types that are backed by a JavaScript value.
type Wrapper interface {
	// JSValue returns a JavaScript value associated with an object.
	JSValue() js.Value
}
