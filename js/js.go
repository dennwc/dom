package js

import "syscall/js"

var (
	global    = js.Global()
	null      = js.Null()
	undefined = js.Undefined()
)

func Get(name string) Value {
	return Value{global.Get(name)}
}

func Set(name string, v Value) {
	global.Set(name, v.JSRef())
}

type Ref = js.Value

type JSRef interface {
	JSRef() Ref
}

func toJS(o interface{}) interface{} {
	switch v := o.(type) {
	case JSRef:
		o = v.JSRef()
	case []Value:
		refs := make([]Ref, 0, len(v))
		for _, ref := range v {
			refs = append(refs, ref.JSRef())
		}
		o = refs
	}
	return o
}

var _ JSRef = Value{}

type Value struct {
	Ref
}

func (v Value) JSRef() Ref {
	return v.Ref
}
func (v Value) String() string {
	if !v.Valid() {
		return ""
	}
	return v.Ref.String()
}
func (v Value) IsNull() bool {
	return v.Ref == null
}
func (v Value) IsUndefined() bool {
	return v.Ref == undefined
}
func (v Value) Valid() bool {
	return !v.IsNull() && !v.IsUndefined()
}
func (v Value) Get(name string) Value {
	return Value{v.Ref.Get(name)}
}
func (v Value) Set(name string, val interface{}) {
	v.Ref.Set(name, toJS(val))
}
func (v Value) Del(name string) {

}
func (v Value) Index(i int) Value {
	return Value{v.Ref.Index(i)}
}
func (v Value) SetIndex(i int, val interface{}) {
	v.Ref.SetIndex(i, toJS(val))
}
func (v Value) Call(name string, args ...interface{}) Value {
	for i, a := range args {
		args[i] = toJS(a)
	}
	return Value{v.Ref.Call(name, args...)}
}
func (v Value) Invoke(args ...interface{}) Value {
	for i, a := range args {
		args[i] = toJS(a)
	}
	return Value{v.Ref.Invoke(args...)}
}
func (v Value) New(args ...interface{}) Value {
	for i, a := range args {
		args[i] = toJS(a)
	}
	return Value{v.Ref.New(args...)}
}
func (v Value) InstanceOf(class Value) bool {
	return v.Ref.InstanceOf(class.Ref)
}
func (v Value) Slice() []Value {
	if !v.Valid() {
		return nil
	}
	n := v.Length()
	vals := make([]Value, 0, n)
	for i := 0; i < n; i++ {
		vals = append(vals, v.Index(i))
	}
	return vals
}

type Callback = js.Callback

func NewCallback(fnc func(v []Value)) Callback {
	return js.NewCallback(func(refs []js.Value) {
		vals := make([]Value, 0, len(refs))
		for _, ref := range refs {
			vals = append(vals, Value{ref})
		}
		fnc(vals)
	})
}

func NewEventCallback(fnc func(v Value)) Callback {
	return NewEventCallbackFlags(0, fnc)
}

func NewEventCallbackFlags(flags int, fnc func(v Value)) Callback {
	return js.NewEventCallback(js.EventCallbackFlag(flags), func(ref js.Value) {
		fnc(Value{ref})
	})
}
