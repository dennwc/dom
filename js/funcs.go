//+build wasm,js

package js

// NewFuncJS creates a function from a JS code string.
//
// Deprecated: use RawFuncOf
//
// Example:
//	 NativeFuncOf("a", "b", "return a+b").Call(a, b)
func NewFuncJS(argsAndCode ...string) Value {
	return NativeFuncOf(argsAndCode...)
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

// NewCallbackGroup creates a new function group on this object.
//
// Deprecated: use NewFuncGroup
func (v Value) NewCallbackGroup() *FuncGroup {
	return v.NewFuncGroup()
}

// NewFuncGroup creates a new function group on this object.
func (v Value) NewFuncGroup() *FuncGroup {
	return &FuncGroup{
		v: v,
	}
}

// CallbackGroup is a list of Go functions attached to an object.
//
// Deprecated: use FuncGroup
type CallbackGroup = FuncGroup

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
