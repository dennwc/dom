package dom

import (
	"fmt"
	"syscall/js"
)

var _ Node = (*Element)(nil)

type Element struct {
	v js.Value
}

func (e *Element) JSValue() js.Value {
	return e.v
}

func (e *Element) AppendChild(n Node) {
	e.v.Call("appendChild", n.JSValue())
}

func (e *Element) SetInnerHTML(s string) {
	e.v.Set("innerHTML", s)
}

func (e *Element) SetAttribute(k string, v interface{}) {
	e.v.Call("setAttribute", k, fmt.Sprint(v))
}

func (e *Element) GetAttribute(k string) js.Value {
	return e.v.Call("getAttribute", k)
}

func (e *Element) Style() *Style {
	return &Style{v: e.v.Get("style")}
}

type Style struct {
	v js.Value
}

func (s *Style) SetWidth(v Unit) {
	s.v.Set("width", v.String())
}

func (s *Style) SetHeight(v Unit) {
	s.v.Set("height", v.String())
}

func (s *Style) SetMarginsRaw(m string) {
	s.v.Set("margin", m)
}

func (s *Style) Set(k string, v interface{}) {
	s.v.Set(k, v)
}
