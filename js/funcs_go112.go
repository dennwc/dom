//+build wasm,go1.12

package js

import "syscall/js"

// NewCallback returns a wrapped callback function.
//
// Invoking the callback in JavaScript will queue the Go function fn for execution.
// This execution happens asynchronously on a special goroutine that handles all callbacks and preserves
// the order in which the callbacks got called.
// As a consequence, if one callback blocks this goroutine, other callbacks will not be processed.
// A blocking callback should therefore explicitly start a new goroutine.
//
// Callback.Release must be called to free up resources when the callback will not be used any more.
func NewCallback(fnc func(v []Value)) Callback {
	return js.NewCallback(func(this js.Value, refs []js.Value) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		fnc(vals)
		return nil
	})
}

// NewCallbackAsync returns a wrapped callback function.
//
// Invoking the callback in JavaScript will queue the Go function fn for execution.
// This execution happens asynchronously.
//
// Callback.Release must be called to free up resources when the callback will not be used any more.
func NewCallbackAsync(fnc func(v []Value)) Callback {
	return js.NewCallback(func(this js.Value, refs []js.Value) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		go fnc(vals)
		return nil
	})
}

// NewFunc returns a wrapped function that will be executed synchronously.
func NewFunc(fnc func(this Value, args []Value) interface{}) Callback {
	return js.NewCallback(func(this js.Value, refs []js.Value) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		v := fnc(Value{this}, vals)
		return ValueOf(v).Ref
	})
}
