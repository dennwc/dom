//+build wasm,go1.12

package js

import "syscall/js"

func unwrap(v Wrapper) Wrapper {
	return v
}

// JSRef is a common interface for object that are backed by a JS object.
//
// Deprecated: see Wrapper
type JSRef interface {
	deprecated()
}

// Wrapper is an alias for syscall/js.Wrapper.
type Wrapper = js.Wrapper
