//+build wasm

package js

import (
	"syscall/js"
)

// Callback is a Go function that got wrapped for use as a JavaScript callback.
type Callback = js.Callback

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
func NewCallbackAsync(fnc func(v []Value)) Callback {
	return js.NewCallback(func(refs []js.Value) {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		go fnc(vals)
	})
}

// NewEventCallback is a shorthand for NewEventCallbackFlags with default flags.
func NewEventCallback(fnc func(v Value)) Callback {
	return NewEventCallbackFlags(0, fnc)
}

// NewEventCallbackFlags returns a wrapped callback function, just like NewCallback, but the callback expects to have
// exactly one argument, the event. Depending on flags, it will synchronously call event.preventDefault,
// event.stopPropagation and/or event.stopImmediatePropagation before queuing the Go function fn for execution.
func NewEventCallbackFlags(flags int, fnc func(v Value)) Callback {
	return js.NewEventCallback(js.EventCallbackFlag(flags), func(ref js.Value) {
		fnc(Value{ref})
	})
}

// NewCallbackGroup creates a new callback group on this object.
func (v Value) NewCallbackGroup() *CallbackGroup {
	return &CallbackGroup{
		v: v,
	}
}

// CallbackGroup is a list of Go callbacks attached to an object.
type CallbackGroup struct {
	v   Value
	cbs []Callback
}

func (g *CallbackGroup) Add(cb Callback) {
	g.cbs = append(g.cbs, cb)
}
func (g *CallbackGroup) Set(name string, fnc func([]Value)) {
	cb := NewCallback(fnc)
	g.v.Set(name, cb)
	g.Add(cb)
}
func (g *CallbackGroup) addEventListener(event string, cb Callback) {
	g.v.Call("addEventListener", event, cb)
}
func (g *CallbackGroup) removeEventListener(event string, cb Callback) {
	g.v.Call("removeEventListener", event, cb)
}
func (g *CallbackGroup) AddEventListener(event string, fnc func(Value)) {
	cb := NewEventCallback(fnc)
	g.addEventListener(event, cb)
	g.Add(cb)
}
func (g *CallbackGroup) ErrorEvent(fnc func(error)) {
	g.AddEventListener("onerror", func(v Value) {
		fnc(Error{v})
	})
}
func (g *CallbackGroup) ErrorEventChan() <-chan error {
	ch := make(chan error, 1)
	g.ErrorEvent(func(err error) {
		select {
		case ch <- err:
		default:
			panic("unhandled error event")
		}
	})
	return ch
}
func (g *CallbackGroup) OneTimeEvent(event string, fnc func(Value)) {
	var cb Callback
	fired := false
	cb = NewEventCallback(func(v Value) {
		if fired {
			panic("one time callback fired twice")
		}
		fnc(v)
		g.removeEventListener(event, cb)
		cb.Release()
	})
	g.addEventListener(event, cb)
	g.Add(cb)
}
func (g *CallbackGroup) OneTimeEventChan(event string) <-chan Value {
	ch := make(chan Value, 1)
	g.OneTimeEvent(event, func(v Value) {
		select {
		case ch <- v:
		default:
			panic("one time callback fired twice")
		}
	})
	return ch
}
func (g *CallbackGroup) OneTimeTrigger(event string) <-chan struct{} {
	ch := make(chan struct{})
	g.OneTimeEvent(event, func(v Value) {
		close(ch)
	})
	return ch
}
func (g *CallbackGroup) Release() {
	for _, cb := range g.cbs {
		cb.Release()
	}
	g.cbs = nil
}

// NewFunction creates a function from JS code string.
//
// Example:
//	 NewFunction("a", "b", "return a+b").Call(a, b)
func NewFunction(argsAndCode ...string) Value {
	args := make([]interface{}, len(argsAndCode))
	for i, v := range argsAndCode {
		args[i] = v
	}
	return New("Function", args...)
}
