//+build !wasm

package js

import "errors"

var (
	global    = Ref{ref: 2}
	null      = Ref{ref: 1}
	undefined = Ref{ref: 0}
)

var errNotImplemented = errors.New("not implemented; build with GOARCH=wasm")

func valueOf(o interface{}) Ref {
	return undefined
}

type Wrapper = interface {
	// JSValue returns a JavaScript value associated with an object.
	JSValue() Ref
}

// Func is a wrapped Go function to be called by JavaScript.
type Func struct {
	Value Ref
}

// Release frees up resources allocated for the function.
// The function must not be invoked after calling Release.
func (f Func) Release() {
	panic(errNotImplemented)
}

func funcOf(fnc func(this Ref, refs []Ref) interface{}) Func {
	panic(errNotImplemented)
}

// Ref is an alias for syscall/js.Value.
type Ref struct {
	ref uint64
}

// Type returns the JavaScript type of the value v. It is similar to JavaScript's typeof operator,
// except that it returns TypeNull instead of TypeObject for null.
func (v Ref) Type() Type {
	panic(errNotImplemented)
}

// Get returns the JavaScript property p of value v.
func (v Ref) Get(k string) Ref {
	return undefined
}

// Set sets the JavaScript property p of value v to ValueOf(x).
func (v Ref) Set(p string, x interface{}) {
	panic(errNotImplemented)
}

// Index returns JavaScript index i of value v.
func (v Ref) Index(i int) Ref {
	panic(errNotImplemented)
}

// SetIndex sets the JavaScript index i of value v to ValueOf(x).
func (v Ref) SetIndex(i int, x interface{}) {
	panic(errNotImplemented)
}

// Length returns the JavaScript property "length" of v.
func (v Ref) Length() int {
	return 0
}

// Call does a JavaScript call to the method m of value v with the given arguments.
// It panics if v has no method m.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Ref) Call(m string, args ...interface{}) Ref {
	panic(errNotImplemented)
}

// Invoke does a JavaScript call of the value v with the given arguments.
// It panics if v is not a function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Ref) Invoke(args ...interface{}) Ref {
	panic(errNotImplemented)
}

// New uses JavaScript's "new" operator with value v as constructor and the given arguments.
// It panics if v is not a function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Ref) New(args ...interface{}) Ref {
	panic(errNotImplemented)
}

// Float returns the value v as a float64. It panics if v is not a JavaScript number.
func (v Ref) Float() float64 {
	panic(errNotImplemented)
}

// Int returns the value v truncated to an int. It panics if v is not a JavaScript number.
func (v Ref) Int() int {
	panic(errNotImplemented)
}

// Bool returns the value v as a bool. It panics if v is not a JavaScript boolean.
func (v Ref) Bool() bool {
	panic(errNotImplemented)
}

// Truthy returns the JavaScript "truthiness" of the value v. In JavaScript,
// false, 0, "", null, undefined, and NaN are "falsy", and everything else is
// "truthy". See https://developer.mozilla.org/en-US/docs/Glossary/Truthy.
func (v Ref) Truthy() bool {
	panic(errNotImplemented)
}

// String returns the value v converted to string according to JavaScript type conversions.
func (v Ref) String() string {
	panic(errNotImplemented)
}

// InstanceOf reports whether v is an instance of type t according to JavaScript's instanceof operator.
func (v Ref) InstanceOf(t Ref) bool {
	panic(errNotImplemented)
}

// Error is an alias for syscall/js.Error.
type Error struct {
	Value Ref
}

// Error implements the error interface.
func (e Error) Error() string {
	return "JavaScript error: undefined"
}

// Type is a type name of a JS value, as returned by "typeof".
type Type int

const (
	TypeObject = Type(iota + 1)
	TypeFunction
)

func typedArrayOf(slice interface{}) Ref {
	panic(errNotImplemented)
}

func releaseTypedArray(v Ref) {
	panic(errNotImplemented)
}
