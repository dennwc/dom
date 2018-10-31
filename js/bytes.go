//+build wasm,js

package js

import (
	"fmt"
	"syscall/js"
)

var _ Wrapper = TypedArray{}

// TypedArray represents a JavaScript typed array.
type TypedArray struct {
	Value
}

// Release frees up resources allocated for the typed array.
// The typed array and its buffer must not be accessed after calling Release.
func (v TypedArray) Release() {
	js.TypedArray{v.Ref}.Release()
}

// TypedArrayOf returns a JavaScript typed array backed by the slice's underlying array.
//
// The supported types are []int8, []int16, []int32, []uint8, []uint16, []uint32, []float32 and []float64.
// Passing an unsupported value causes a panic.
//
// TypedArray.Release must be called to free up resources when the typed array will not be used any more.
func TypedArrayOf(o interface{}) TypedArray {
	v := js.TypedArrayOf(toJS(o))
	return TypedArray{Value{v.Value}}
}

var _ Wrapper = (*Memory)(nil)

type Memory struct {
	v TypedArray
	p []byte
}

func (m *Memory) Bytes() []byte {
	return m.p
}

// CopyFrom copies binary data from JS object into Go buffer.
func (m *Memory) CopyFrom(v Wrapper) error {
	var src Value
	switch v := v.(type) {
	case Value:
		src = v
	case TypedArray:
		src = v.Value
	default:
		src = Value{v.JSValue()}
	}

	switch {
	case src.InstanceOfClass("Uint8Array"):
		m.v.Call("set", src)
		return nil
	case src.InstanceOfClass("Blob"):
		r := New("FileReader")

		cg := r.NewCallbackGroup()
		defer cg.Release()
		done := cg.OneTimeTrigger("loadend")
		errc := cg.ErrorEventChan()

		r.Call("readAsArrayBuffer", src)
		select {
		case err := <-errc:
			return err
		case <-done:
		}
		cg.Release()
		arr := New("Uint8Array", r.Get("result"))
		return m.CopyFrom(arr)
	default:
		return fmt.Errorf("unsupported source type")
	}
}

func (m *Memory) JSValue() Ref {
	return m.v.JSValue()
}

func (m *Memory) Release() {
	m.v.Release()
}

// MMap exposes memory of p to JS.
//
// Release must be called to free up resources when the memory will not be used any more.
func MMap(p []byte) *Memory {
	v := TypedArrayOf(p)
	return &Memory{p: p, v: v}
}
