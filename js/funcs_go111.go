//+build wasm,!go1.12

package js

import "syscall/js"

// Callback is a wrapped Go function to be called by JavaScript.
type Callback = js.Callback

// Func is a wrapped Go function to be called by JavaScript.
type Func = Callback

// NewCallback returns a wrapped callback function.
//
// Invoking the callback in JavaScript will queue the Go function fn for execution.
// This execution happens asynchronously on a special goroutine that handles all callbacks and preserves
// the order in which the callbacks got called.
// As a consequence, if one callback blocks this goroutine, other callbacks will not be processed.
// A blocking callback should therefore explicitly start a new goroutine.
//
// Callback.Release must be called to free up resources when the callback will not be used any more.
//
// Deprecated: use CallbackOf
func NewCallback(fnc func(v []Value)) Callback {
	return CallbackOf(fnc)
}

// CallbackOf returns a wrapped callback function.
//
// Invoking the callback in JavaScript will queue the Go function fn for execution.
// This execution happens asynchronously on a special goroutine that handles all callbacks and preserves
// the order in which the callbacks got called.
// As a consequence, if one callback blocks this goroutine, other callbacks will not be processed.
// A blocking callback should therefore explicitly start a new goroutine.
//
// Callback.Release must be called to free up resources when the callback will not be used any more.
func CallbackOf(fnc func(v []Value)) Callback {
	return js.NewCallback(func(refs []js.Value) {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		fnc(vals)
	})
}

// NewCallbackAsync returns a wrapped callback function.
//
// Invoking the callback in JavaScript will queue the Go function fn for execution.
// This execution happens asynchronously.
//
// Callback.Release must be called to free up resources when the callback will not be used any more.
//
// Deprecated: use AsyncCallbackOf
func NewCallbackAsync(fnc func(v []Value)) Func {
	return AsyncCallbackOf(fnc)
}

// AsyncCallbackOf returns a wrapped callback function.
//
// Invoking the callback in JavaScript will queue the Go function fn for execution.
// This execution happens asynchronously.
//
// Callback.Release must be called to free up resources when the callback will not be used any more.
func AsyncCallbackOf(fnc func(v []Value)) Func {
	return js.NewCallback(func(refs []js.Value) {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		go fnc(vals)
	})
}
