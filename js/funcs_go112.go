//+build wasm,go1.12

package js

import "syscall/js"

// Callback is a wrapped Go function to be called by JavaScript.
//
// Deprecated: use Func
type Callback = Func

// Func is a wrapped Go function to be called by JavaScript.
type Func = js.Func

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
func NewCallback(fnc func(v []Value)) Func {
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
func CallbackOf(fnc func(v []Value)) Func {
	return js.FuncOf(func(this js.Value, refs []js.Value) interface{} {
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
//
// Deprecated: AsyncCallbackOf
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
	return js.FuncOf(func(this js.Value, refs []js.Value) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		go fnc(vals)
		return nil
	})
}

// NewFunc returns a wrapped function that will be executed synchronously.
//
// Deprecated: use FuncOf
func NewFunc(fnc func(this Value, args []Value) interface{}) Func {
	return FuncOf(fnc)
}

// FuncOf returns a wrapped function.
//
// Invoking the JavaScript function will synchronously call the Go function fn with the value of JavaScript's
// "this" keyword and the arguments of the invocation.
// The return value of the invocation is the result of the Go function mapped back to JavaScript according to ValueOf.
//
// A wrapped function triggered during a call from Go to JavaScript gets executed on the same goroutine.
// A wrapped function triggered by JavaScript's event loop gets executed on an extra goroutine.
// Blocking operations in the wrapped function will block the event loop.
// As a consequence, if one wrapped function blocks, other wrapped funcs will not be processed.
// A blocking function should therefore explicitly start a new goroutine.
//
// Func.Release must be called to free up resources when the function will not be used any more.
func FuncOf(fnc func(this Value, args []Value) interface{}) Func {
	return js.FuncOf(func(this js.Value, refs []js.Value) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		v := fnc(Value{this}, vals)
		return ValueOf(v).Ref
	})
}
