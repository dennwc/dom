//+build wasm

// Package JS provides additional functionality on top of syscall/js package for WASM.
package js

import (
	"sync"
	"syscall/js"
)

var (
	global    = js.Global()
	null      = js.Null()
	undefined = js.Undefined()

	object = global.Get("Object")
	array  = global.Get("Array")
)

var (
	mu      sync.RWMutex
	classes = make(map[string]Value)
)

// Obj is an alias for map[string]interface{}.
type Obj = map[string]interface{}

// Arr is an alias for []interface{}.
type Arr = []interface{}

// Get is a shorthand for Global().Get().
func Get(name string, path ...string) Value {
	return Value{global}.Get(name, path...)
}

// Set is a shorthand for Global().Set().
func Set(name string, v interface{}) {
	global.Set(name, toJS(v))
}

// Call is a shorthand for Global().Call().
func Call(name string, args ...interface{}) Value {
	return Value{global}.Call(name, args...)
}

// Class searches for a class in global scope.
// It caches results, so the lookup should be faster than calling Get.
func Class(class string) Value {
	switch class {
	case "Object":
		return Value{object}
	case "Array":
		return Value{array}
	}
	mu.RLock()
	v := classes[class]
	mu.RUnlock()
	if v.isZero() {
		v = Get(class)
		mu.Lock()
		classes[class] = v
		mu.Unlock()
	}
	return v
}

// New searches for a class in global scope and creates a new instance of that class.
func New(class string, args ...interface{}) Value {
	v := Class(class)
	return v.New(args...)
}

// Ref is an alias for syscall/js.Value.
type Ref = js.Value

// JSRef is a common interface for object that are backed by a JS object.
//
// Deprecated: see Wrapper
type JSRef interface {
	deprecated()
}

// Wrapper is an alias for syscall/js.Wrapper.
type Wrapper = js.Wrapper

// Error is an alias for syscall/js.Error.
type Error = js.Error

// Type is a type name of a JS value, as returned by "typeof".
type Type = js.Type

const (
	TypeObject   = js.TypeObject
	TypeFunction = js.TypeFunction
)

// Type is an analog for JS "typeof" operator.
func (v Value) Type() Type {
	return v.Ref.Type()
}

// Object returns an Object JS class.
func Object() Value {
	return Value{object}
}

// Array returns an Array JS class.
func Array() Value {
	return Value{array}
}

// NewObject creates an empty JS object.
func NewObject() Value {
	return Object().New()
}

// NewArray creates an empty JS array.
func NewArray() Value {
	return Array().New()
}

func toJS(o interface{}) interface{} {
	switch v := o.(type) {
	case Wrapper:
		// pass directly
	case []Value:
		refs := make([]interface{}, 0, len(v))
		for _, ref := range v {
			refs = append(refs, ref.JSValue())
		}
		o = refs
	case []js.Value:
		refs := make([]interface{}, 0, len(v))
		for _, ref := range v {
			refs = append(refs, ref)
		}
		o = refs
	}
	return o
}

var _ Wrapper = Value{}

// Value is a convenience wrapper for syscall/js.Value.
// It provides some additional functionality, while storing no additional state.
// Its safe to instantiate Value directly, by wrapping syscall/js.Value.
type Value struct {
	Ref
}

func (v Value) isZero() bool {
	return v == (Value{})
}

// JSRef returns a JS object reference as defined by syscall/js.
//
// Deprecated: use JSValue
func (v Value) JSRef() Ref {
	return v.JSValue()
}

// JSValue implements Wrapper interface.
func (v Value) JSValue() Ref {
	return v.Ref
}

// String converts a value to a string.
func (v Value) String() string {
	if !v.Valid() {
		return ""
	}
	return v.Ref.String()
}

// IsNull checks if a value represents JS null object.
func (v Value) IsNull() bool {
	return v.Ref == null
}

// IsUndefined checks if a value represents JS undefined object.
func (v Value) IsUndefined() bool {
	return v.Ref == undefined
}

// Valid checks if object is defined and not null.
func (v Value) Valid() bool {
	return !v.isZero() && !v.IsNull() && !v.IsUndefined()
}

// Get returns the JS property by name.
func (v Value) Get(name string, path ...string) Value {
	ref := v.Ref.Get(name)
	for _, p := range path {
		ref = ref.Get(p)
	}
	return Value{ref}
}

// Set sets the JS property to ValueOf(x).
func (v Value) Set(name string, val interface{}) {
	v.Ref.Set(name, ValueOf(val))
}

// TODO: Del

// Index returns JS index i of value v.
func (v Value) Index(i int) Value {
	return Value{v.Ref.Index(i)}
}

// SetIndex sets the JavaScript index i of value v to ValueOf(x).
func (v Value) SetIndex(i int, val interface{}) {
	v.Ref.SetIndex(i, ValueOf(val))
}

// Call does a JavaScript call to the method m of value v with the given arguments.
// It panics if v has no method m.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) Call(name string, args ...interface{}) Value {
	for i, a := range args {
		args[i] = ValueOf(a)
	}
	return Value{v.Ref.Call(name, args...)}
}

// Invoke does a JavaScript call of the value v with the given arguments.
// It panics if v is not a function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) Invoke(args ...interface{}) Value {
	for i, a := range args {
		args[i] = ValueOf(a)
	}
	return Value{v.Ref.Invoke(args...)}
}

// New uses JavaScript's "new" operator with value v as constructor and the given arguments.
// It panics if v is not a function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) New(args ...interface{}) Value {
	for i, a := range args {
		args[i] = ValueOf(a)
	}
	return Value{v.Ref.New(args...)}
}

// InstanceOf reports whether v is an instance of type t according to JavaScript's instanceof operator.
func (v Value) InstanceOf(class Wrapper) bool {
	return v.Ref.InstanceOf(class.JSValue())
}

// InstanceOfClass reports whether v is an instance of named type according to JavaScript's instanceof operator.
func (v Value) InstanceOfClass(class string) bool {
	return v.InstanceOf(Class(class))
}

// Slice converts JS Array to a Go slice of JS values.
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

// ValueOf returns x as a JavaScript value:
//
//  | Go                     | JavaScript             |
//  | ---------------------- | ---------------------- |
//  | js.Value               | [its value]            |
//  | js.TypedArray          | typed array            |
//  | js.Callback            | function               |
//  | nil                    | null                   |
//  | bool                   | boolean                |
//  | integers and floats    | number                 |
//  | string                 | string                 |
//  | []interface{}          | new array              |
//  | map[string]interface{} | new object             |
func ValueOf(o interface{}) Value {
	return Value{js.ValueOf(toJS(o))}
}
