package js

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
	return funcOf(func(this Ref, refs []Ref) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		fnc(vals)
		return nil
	})
}

// AsyncCallbackOf returns a wrapped callback function.
//
// Invoking the callback in JavaScript will queue the Go function fn for execution.
// This execution happens asynchronously.
//
// Callback.Release must be called to free up resources when the callback will not be used any more.
func AsyncCallbackOf(fnc func(v []Value)) Func {
	return funcOf(func(this Ref, refs []Ref) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		go fnc(vals)
		return nil
	})
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
	return funcOf(func(this Ref, refs []Ref) interface{} {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		v := fnc(Value{this}, vals)
		return ValueOf(v).Ref
	})
}

// NativeFuncOf creates a function from a JS code string.
//
// Example:
//	 NativeFuncOf("a", "b", "return a+b").Call(a, b)
func NativeFuncOf(argsAndCode ...string) Value {
	args := make([]interface{}, len(argsAndCode))
	for i, v := range argsAndCode {
		args[i] = v
	}
	return New("Function", args...)
}

// NewEventCallback is a shorthand for NewEventCallbackFlags with default flags.
func NewEventCallback(fnc func(v Value)) Func {
	return CallbackOf(func(v []Value) {
		fnc(v[0])
	})
}

// NewFuncGroup creates a new function group on this object.
func (v Value) NewFuncGroup() *FuncGroup {
	return &FuncGroup{
		v: v,
	}
}

// FuncGroup is a list of Go functions attached to an object.
type FuncGroup struct {
	v     Value
	funcs []Func
}

func (g *FuncGroup) Add(cb Func) {
	g.funcs = append(g.funcs, cb)
}
func (g *FuncGroup) Set(name string, fnc func([]Value)) {
	cb := CallbackOf(fnc)
	g.v.Set(name, cb)
	g.Add(cb)
}
func (g *FuncGroup) addEventListener(event string, cb Func) {
	g.v.Call("addEventListener", event, cb)
}
func (g *FuncGroup) removeEventListener(event string, cb Func) {
	g.v.Call("removeEventListener", event, cb)
}
func (g *FuncGroup) AddEventListener(event string, fnc func(Value)) {
	cb := NewEventCallback(fnc)
	g.addEventListener(event, cb)
	g.Add(cb)
}
func (g *FuncGroup) ErrorEvent(fnc func(error)) {
	g.AddEventListener("onerror", func(v Value) {
		fnc(Error{v.Ref})
	})
}
func (g *FuncGroup) ErrorEventChan() <-chan error {
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
func (g *FuncGroup) OneTimeEvent(event string, fnc func(Value)) {
	var cb Func
	fired := false
	cb = NewEventCallback(func(v Value) {
		if fired {
			panic("one time callback fired twice")
		}
		fired = true
		fnc(v)
		g.removeEventListener(event, cb)
		cb.Release()
	})
	g.addEventListener(event, cb)
	g.Add(cb)
}
func (g *FuncGroup) OneTimeEventChan(event string) <-chan Value {
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
func (g *FuncGroup) OneTimeTrigger(event string) <-chan struct{} {
	ch := make(chan struct{})
	g.OneTimeEvent(event, func(v Value) {
		close(ch)
	})
	return ch
}
func (g *FuncGroup) Release() {
	for _, f := range g.funcs {
		f.Release()
	}
	g.funcs = nil
}
