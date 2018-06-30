package dom

import (
	"fmt"
	"syscall/js"
)

var _ Node = (*Element)(nil)

func AsElement(v js.Value) *Element {
	if !IsValid(v) {
		return nil
	}
	return &Element{NodeBase{v}}
}

func AsNodeList(v js.Value) NodeList {
	if !IsValid(v) {
		return nil
	}
	arr := make(NodeList, v.Length())
	for i := range arr {
		arr[i] = AsElement(v.Index(i))
	}
	return arr
}

var _ Node = (*Element)(nil)

type Element struct {
	NodeBase
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
